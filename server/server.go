// Package server implements a simple server to return concordance (word count + sentence location) for sample text.
package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/composer22/clidemo/logger"
)

// Server is the main structure that represents a server instance.
type Server struct {
	mu      sync.Mutex            // For locking access to server params.
	info    *Info                 // Basic server information.
	opts    *Options              // Original options and info for creating the server.
	running bool                  // Is the server running?
	log     *logger.Logger        // Log instance for recording error and other messages.
	jobq    chan *parseJob        // Channel to send jobs.
	mw      *middlewarePreprocess // Handler for all incomping routes
	wg      sync.WaitGroup        // Synchronize close() of job channel.
	stats   *Status               // Server statistics since it started.
}

// Info provides basic information about the running server.
type Info struct {
	Version    string `json:"version"`        // Version of the server.
	Name       string `json:"name"`           // The name of the server.
	UUID       string `json:"UUID"`           // Unique ID of the server.
	Port       int    `json:"port"`           // Port the server is listening on.
	MaxConn    int    `json:"maxConnections"` // The maximum concurrent connections accepted.
	MaxWorkers int    `json:"maxWorkers"`     // The maximum numer of workers allowed to run.
	Debug      bool   `json:"debugEnabled"`   // Is debugging enabled on the server.

}

// stats contains runtime statistics.
type Status struct {
	Start        time.Time        `json:"startTime"`    // The start time of the server.
	RequestCount int64            `json:"requestCount"` // How many requests came in to the server.
	InBytes      int64            `json:"inBytes"`      // Size of the requests in bytes.
	PageRequests map[string]int64 `json:"pageRequests"` // How many requests came into each page.
}

// New is a factory function that returns a new server instance.
func New(opts *Options) *Server {
	log := logger.New(-1, false)

	// Server information.
	info := &Info{
		Version:    version,
		Name:       opts.Name,
		UUID:       createV4UUID(),
		Port:       opts.Port,
		MaxConn:    opts.MaxConn,
		MaxWorkers: opts.MaxWorkers,
		Debug:      opts.Debug,
	}

	// Stat information.
	st := &Status{
		Start:        time.Now(),
		PageRequests: make(map[string]int64),
	}

	// Construct server.
	s := &Server{
		info:    info,
		opts:    opts,
		jobq:    make(chan *parseJob),
		log:     log,
		stats:   st,
		running: false,
	}

	if s.info.Debug {
		s.log.SetLogLevel(logger.Debug)
	}

	// Setup the routes and middleware.
	mux := http.NewServeMux()
	mux.HandleFunc(httpRouteAliveV1, s.aliveHandler)
	mux.HandleFunc(httpRouteParseV1, s.parseHandler)
	mux.HandleFunc(httpRouteStatusV1, s.statusHandler)
	s.mw = &middlewarePreprocess{serv: s, handler: mux}

	// Trap signals
	s.handleSignals()
	return s
}

// PrintVersionAndExit prints the version of the server then exits.
func PrintVersionAndExit() {
	fmt.Printf("clidemo version %s\n", version)
	os.Exit(0)
}

// Start spins up the server to accept incoming connections.
func (s *Server) Start() {
	s.log.Infof("Starting clidemo version %s\n", version)
	s.mu.Lock()
	s.stats.Start = time.Now()
	s.running = true

	// Spin off the worker processes.
	for i := 0; i < s.info.MaxWorkers; i++ {
		s.wg.Add(1)
		go parseWorker(s.jobq, &s.wg)
	}

	s.mu.Unlock()
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.info.Port), s.mw)
	if err != nil {
		s.log.Emergencyf("%s\n", err)
	}
}

// handleSignals responds to operating system interrupts such as application kills.
func (s *Server) handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			s.log.Infof("Server received signal: %v\n", sig)
			s.log.Infof("Stopping all workers...")
			close(s.jobq)
			s.wg.Wait()
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
		http.Error(w, invalidBody, http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(b, &data); err != nil {
		http.Error(w, invalidJSONText, http.StatusBadRequest)
		return
	}
	d, e := data["text"].(string)
	if e {
		http.Error(w, invalidJSONAttribute, http.StatusBadRequest)
	}

	// Send a parse request to a parse worker and wait for it to complete.
	job := parseJob{
		Source: d,
		DoneCh: make(chan bool),
	}
	s.jobq <- &job
	<-job.DoneCh

	w.Write([]byte(job.ResultJSON))
}

// statusHandler handles a client request for server information and statistics.
func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	if s.invalidMethod(w, r, httpGet) {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	inf, _ := json.Marshal(s.info)
	op, _ := json.Marshal(s.opts)
	st, _ := json.Marshal(s.stats)
	w.Write([]byte(fmt.Sprintf(`{"info":%s,"options":%s,"stats":%s}`, inf, op, st)))
}

// incrementStats increments the statistics for the request being handled by the server.
func (s *Server) incrementStats(r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.RequestCount++
	s.stats.PageRequests[r.URL.Path]++
	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		s.stats.InBytes += int64(len(b))
	}
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
		http.Error(w, invalidMediaType, http.StatusUnsupportedMediaType)
		return true
	}
	return false
}

// invalidMethod validates that the http method is acceptable for processing this route.
func (s *Server) invalidMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		http.Error(w, invalidMediaType, http.StatusUnsupportedMediaType)
		return true
	}
	return false
}

// isRunning returns a boolean representing whether the server is running or not.
func (s *Server) isRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
