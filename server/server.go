// Package server implements a simple server to return concordance (word count + sentence location) for sample text.
package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

// Server is the main structure that represents a server instance.
type Server struct {
	mu       sync.Mutex     // For locking access to server params.
	info     Info           // Basic server information.
	infoJSON []byte         // Basic server information as json.
	opts     *Options       // Original options and info for creating the server.
	running  bool           // Is the server running?
	jobq     chan *parseJob // Channel to send jobs.
	wg       sync.WaitGroup // Synchronize close() of job channel.
	stats    *Status        // Server statistics since it started.
}

// Info provides basic information about the server.
type Info struct {
	UUID    string `json:"UUID"`    // Unique ID of the server.
	Version string `json:"version"` // Version of the server.
	Port    int    `json:"port"`    // Port the server is listening on.
}

// stats contains runtime statistics.
type Status struct {
	Start        time.Time `json:"startTime"`    // The start time of the server.
	RequestCount int64     `json:"requestCount"` // How many requests came in to the server.
	InBytes      int64     `json:"inBytes"`      // Size of the requests in bytes.
}

// New is a factory function that returns a new server instance.
func New(opts *Options) *Server {

	// Server information.
	info := Info{
		UUID:    createV4UUID(),
		Version: version,
		Port:    opts.Port,
	}

	// Create json version of the info.
	b, err := json.Marshal(info)
	if err != nil {
		log.Fatalf("[FATAL] Error marshalling info json: %+v", err)
	}

	// Stat information.
	st := &Status{
		Start: time.Now(),
	}

	// Construct server.
	s := &Server{
		info:     info,
		infoJSON: b,
		opts:     opts,
		jobq:     make(chan *parseJob),
		stats:    st,
		running:  false,
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
	log.Printf("[INFO] Starting clidemo version %s\n", version)
	s.mu.Lock()
	s.stats.Start = time.Now()
	s.running = true
	s.mu.Unlock()

	// Spin off the worker processes.
	for i := 0; i < s.opts.MaxConn; i++ {
		s.wg.Add(1)
		go parseWorker(s.jobq, &s.wg)
	}

	// Setup the routes and middleware, and serve.
	mux := http.NewServeMux()
	mux.HandleFunc(httpRouteAliveV1, s.aliveHandler)
	mux.HandleFunc(httpRouteParseV1, s.parseHandler)
	mux.HandleFunc(httpRouteStatusV1, s.statusHandler)
	mw := &middlewarePreprocess{serv: s, handler: mux}
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.opts.Port), mw)
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}
}

// handleSignals responds to operating system interrupts such as application kills.
func (s *Server) handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			s.mu.Lock()
			if s.opts.Debug {
				log.Printf("[DEBUG] Trapped signal: %v\n", sig)
			}
			s.mu.Unlock()
			log.Println("[INFO] Server closing all jobs...")
			close(s.jobq)
			s.wg.Wait()
			log.Println("[INFO] Server exiting...")
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
	if s.invalidMethod(w, r, httpGet) {
		return
	}

	// Read the json in for the request.
	var data map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
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

	// Send a parse request to a parse worker.
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
	st, _ := json.Marshal(s.stats)
	w.Write([]byte(fmt.Sprintf(`{"info":%s,"stats":%s}`, s.infoJSON, st)))
}

// incrementStats increments the statistics for the request being handled by the server.
func (s *Server) incrementStats(r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.RequestCount++
	cl, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err == nil {
		s.stats.InBytes += int64(cl)
	}
}

// initResponseHeader sets up the common http response headers for the return of all json calls.
func (s *Server) initResponseHeader(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.Header().Add("Date", time.Now().UTC().Format(time.RFC1123Z))
	if s.opts.Name != "" {
		w.Header().Add("Server", s.opts.Name)
	}
	w.Header().Add("X-Request-ID", createV4UUID())
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
