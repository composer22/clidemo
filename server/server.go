// Package server implements a simple server to return concordance (word count + sentence location) for sample text.
package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"

	// Allow dynamic profiling.
	_ "net/http/pprof"

	"github.com/composer22/clidemo/auth"
	"github.com/composer22/clidemo/logger"
)

type requestLogEntry struct {
	Method        string      `json:"method"`
	URL           *url.URL    `json:"url"`
	Proto         string      `json:"proto"`
	Header        http.Header `json:"header"`
	Body          string      `json:"body"`
	ContentLength int64       `json:"contentLength"`
	Host          string      `json:"host"`
	RemoteAddr    string      `json:"remoteAddr"`
	RequestURI    string      `json:"requestURI"`
	Trailer       http.Header `json:"trailer"`
}

// Server is the main structure that represents a server instance.
type Server struct {
	mu       sync.Mutex         // For locking access to server params.
	info     *Info              // Basic server information.
	opts     *Options           // Original options and info for creating the server.
	running  bool               // Is the server running?
	auth     *auth.Auth         // Authorization lookup
	log      *logger.Logger     // Log instance for recording error and other messages.
	jobq     chan *parseJob     // Channel to send jobs.
	srvr     *http.Server       // HTTP server.
	listener *ThrottledListener // Optional listener for connections.
	wg       sync.WaitGroup     // Synchronize close() of job channel.
	stats    *Status            // Server statistics since it started.
}

// New is a factory function that returns a new server instance.
func New(opts *Options, addedOptions ...func(*Server)) *Server {
	s := &Server{
		info: InfoNew(func(i *Info) {
			i.Name = opts.Name
			i.Hostname = opts.Hostname
			i.Port = opts.Port
			i.ProfPort = opts.ProfPort
			i.MaxConn = opts.MaxConn
			i.MaxWorkers = opts.MaxWorkers
			i.Debug = opts.Debug
		}),
		opts:    opts,
		auth:    auth.New(),
		log:     logger.New(logger.UseDefault, false),
		stats:   StatusNew(),
		running: false,
	}

	if s.info.Debug {
		s.log.SetLogLevel(logger.Debug)
	}

	// Setup the routes, middleware, and server.
	mux := http.NewServeMux()
	mux.HandleFunc(httpRouteAliveV1, s.aliveHandler)
	mux.HandleFunc(httpRouteParseV1, s.parseHandler)
	mux.HandleFunc(httpRouteStatusV1, s.statusHandler)
	s.srvr = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.info.Hostname, s.info.Port),
		Handler:      &Middleware{serv: s, handler: mux},
		ReadTimeout:  TCPReadTimeout,
		WriteTimeout: TCPWriteTimeout,
	}

	s.handleSignals() // Evoke trap signals handler

	// Additional hook for specialized custom options
	for _, option := range addedOptions {
		option(s)
	}
	return s
}

// PrintVersionAndExit prints the version of the server then exits.
func PrintVersionAndExit() {
	fmt.Printf("clidemo version %s\n", version)
	os.Exit(0)
}

// Start spins up the server to accept incoming connections.
func (s *Server) Start() {
	var err error

	s.log.Infof("Starting clidemo version %s\n", version)
	s.mu.Lock()

	s.listener, err = ThrottledListenerNew(s.srvr.Addr, s.info.MaxConn)
	if err != nil {
		s.mu.Unlock()
		s.log.Emergencyf("%s\n", err)
	}

	// Spin off the worker processes.
	s.jobq = make(chan *parseJob)
	for i := 0; i < s.info.MaxWorkers; i++ {
		s.wg.Add(1)
		go parseWorker(s.jobq, &s.wg)
	}

	// Pprof http endpoint for the profiler.
	if s.info.ProfPort > 0 {
		s.StartProfiler()
	}

	s.stats.Start = time.Now()
	s.running = true
	s.mu.Unlock()
	s.srvr.Serve(s.listener)
}

// StartProfiler is called to enable dynamic profiling.
func (s *Server) StartProfiler() {
	s.log.Infof("Starting profiling on http port %d", s.opts.ProfPort)
	hp := fmt.Sprintf("%s:%d", s.info.Hostname, s.info.ProfPort)
	go func() {
		err := http.ListenAndServe(hp, nil)
		if err != nil {
			s.log.Emergencyf("Error starting profile monitoring service: %s", err)
		}
	}()
}

// Shutdown takes down the server gracefully back to an initialize state.
func (s *Server) Shutdown() bool {
	if !s.isRunning() {
		return true
	}

	s.log.Infof("BEGIN server service stop.")

	s.mu.Lock()

	s.log.Infof("\tStopping server listener...")
	s.listener.Stop()
	s.srvr.SetKeepAlivesEnabled(false)

	s.log.Infof("\tStopping workers...")
	close(s.jobq)
	s.wg.Wait()

	s.running = false
	s.jobq = nil
	s.listener = nil
	s.mu.Unlock()

	// Sleep a bit to allow all connections to be released
	var maxTimeout time.Duration
	if TCPReadTimeout > TCPWriteTimeout {
		maxTimeout = TCPReadTimeout
	} else {
		maxTimeout = TCPWriteTimeout
	}
	maxTimeout = maxTimeout + (1 * time.Second)
	s.log.Infof("\tAllowing all connections to close for %s...", maxTimeout.String())
	time.Sleep(maxTimeout)

	s.log.Infof("END server service stop.")
	return true
}

// handleSignals responds to operating system interrupts such as application kills.
func (s *Server) handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			s.log.Infof("Server received signal: %v\n", sig)
			s.Shutdown()
			s.log.Infof("Server exiting.")
			os.Exit(0)
		}
	}()
}

// aliveHandler handles a client "is the server alive" request.
func (s *Server) aliveHandler(w http.ResponseWriter, r *http.Request) {
	if s.invalidMethod(w, r, httpGet) {
		return
	}
}

// parseHandler handles a parse request from the client and returns a json result.
func (s *Server) parseHandler(w http.ResponseWriter, r *http.Request) {
	if s.invalidMethod(w, r, httpPost) {
		return
	}

	// Read the json in for the request.
	var data map[string]interface{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, InvalidBody, http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(b, &data); err != nil {
		http.Error(w, InvalidJSONText, http.StatusBadRequest)
		return
	}
	d, ok := data["text"].(string)
	if !ok {
		http.Error(w, InvalidJSONAttribute, http.StatusBadRequest)
		return
	}

	// Send a parse request to a parse worker and wait for it to complete.
	job := parseJob{
		Source: d,
		DoneCh: make(chan bool),
	}
	s.jobq <- &job
	<-job.DoneCh
	w.Write([]byte(fmt.Sprintf(`{"result":%s}`, job.Result)))
}

// statusHandler handles a client request for server information and statistics.
func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	if s.invalidMethod(w, r, httpGet) {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener != nil {
		s.stats.ConnNumAvail = s.listener.GetConnNumAvail() // Get latest live connection count.
	}
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)
	b, _ := json.Marshal(
		&struct {
			Info    *Info             `json:"info"`
			Options *Options          `json:"options"`
			Stats   *Status           `json:"stats"`
			Memory  *runtime.MemStats `json:"memStats"`
		}{
			Info:    s.info,
			Options: s.opts,
			Stats:   s.stats,
			Memory:  mStats,
		})
	w.Write(b)
}

// incrementStats increments the statistics for the request being handled by the server.
func (s *Server) incrementStats(r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.IncrRequestStats(r.ContentLength)
	s.stats.IncrRouteStats(r.URL.Path, r.ContentLength)
}

// initResponseHeader sets up the common http response headers for the return of all json calls.
func (s *Server) initResponseHeader(w http.ResponseWriter) {
	header := w.Header()
	header.Add("Content-Type", "application/json;charset=utf-8")
	header.Add("Date", time.Now().UTC().Format(time.RFC1123Z))
	if s.info.Name != "" {
		header.Add("Server", s.info.Name)
	}
	header.Add("X-Request-ID", createV4UUID())
}

// invalidHeader validates that the header information is acceptable for processing the request from the client.
func (s *Server) invalidHeader(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Content-Type") != "application/json" || r.Header.Get("Accept") != "application/json" {
		http.Error(w, InvalidMediaType, http.StatusUnsupportedMediaType)
		return true
	}
	return false
}

// invalidAuth validates that the Authorization token is valid for using the API
func (s *Server) invalidAuth(w http.ResponseWriter, r *http.Request) bool {
	if !s.auth.Valid(strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)) {
		http.Error(w, InvalidAuthorization, http.StatusUnauthorized)
		return true
	}
	return false
}

// invalidMethod validates that the http method is acceptable for processing this route.
func (s *Server) invalidMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		http.Error(w, InvalidMethod, http.StatusMethodNotAllowed)
		return true
	}
	return false
}

// LogRequest logs the http request information into the logger.
func (s *Server) LogRequest(r *http.Request) {
	var cl int64
	if r.ContentLength > 0 {
		cl = r.ContentLength
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		body = []byte("Could not parse body")
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // We need to set the body back after we read it.

	b, _ := json.Marshal(&requestLogEntry{
		Method:        r.Method,
		URL:           r.URL,
		Proto:         r.Proto,
		Header:        r.Header,
		Body:          string(body),
		ContentLength: cl,
		Host:          r.Host,
		RemoteAddr:    r.RemoteAddr,
		RequestURI:    r.RequestURI,
		Trailer:       r.Trailer,
	})
	s.log.Infof(`{"request":%s}`, string(b))
}

// isRunning returns a boolean representing whether the server is running or not.
func (s *Server) isRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
