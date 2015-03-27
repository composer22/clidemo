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
	start    time.Time      // The start time of the server.
	stats                   // Server statistics since it started.
}

// Info provides basic information about the server.
type Info struct {
	UUID    string `json:"UUID"`    // Unique ID of the server.
	Version string `json:"version"` // Version of the server.
	Port    int    `json:"port"`    // Port the server is listening on.
}

// stats contains runtime statistics.
type stats struct {
	requestCount int64 // How many requests came in to the server.
	inBytes      int64 // Size of the requests in bytes.
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

	// Construct server.
	s := &Server{
		info:     info,
		infoJSON: []byte(fmt.Sprintf("{\"info\":%s}", b)),
		opts:     opts,
		jobq:     make(chan *parseJob),
		running:  false,
		start:    time.Now(),
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
	s.start = time.Now()
	s.running = true
	s.mu.Unlock()

	for i := 0; i < s.opts.MaxConn; i++ {
		s.wg.Add(1)
		go parseWorker(s.jobq, &s.wg)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(httpRouteAliveV1, s.aliveHandler)
	mux.HandleFunc(httpRouteParseV1, s.parseHandler)
	mux.HandleFunc(httpRouteStatusV1, s.statusHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.opts.Port), mux)
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

	// Validate request header.
	if s.invalidHeader(w, r, httpGet) {
		return
	}

	s.initResponseHeader(w)
}

// parseHandler handles a parse request from the client and returns a json result.
func (s *Server) parseHandler(w http.ResponseWriter, r *http.Request) {

	// Validate request header.
	if s.invalidHeader(w, r, httpGet) {
		return
	}

	// Read the json in for the request.
	var data map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.errorHandler(w, r, invalidBody, http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(b, &data); err != nil {
		s.errorHandler(w, r, invalidJSONText, http.StatusBadRequest)
		return
	}
	d, e := data["text"].(string)
	if e {
		s.errorHandler(w, r, invalidJSONAttribute, http.StatusBadRequest)
	}

	// Send a parse request to a parse worker.
	job := parseJob{
		Source: d,
		DoneCh: make(chan bool),
	}
	s.jobq <- &job
	<-job.DoneCh

	s.initResponseHeader(w)
	w.Write([]byte(job.ResultJSON))
}

// statusHandler handles a client request for server status information.
func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {

	// Validate request header.
	if s.invalidHeader(w, r, httpGet) {
		return
	}

	// TODO Handle statistics and information
	s.initResponseHeader(w)
	// TODO w.Write()
}

// errorHandler wraps a standard response for any invalid condition found by the other http handlers.
func (s *Server) errorHandler(w http.ResponseWriter, r *http.Request, message string, status int) {
	s.initResponseHeader(w)
	w.WriteHeader(status)
	if message != "" {
		fmt.Fprintf(w, `{"error":"%s"}`, message)
	}
}

// invalidHeader validates that the header information is acceptable for processing the request from the client.
func (s *Server) invalidHeader(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method ||
		r.Header.Get("Content-Type") != "application/json" ||
		r.Header.Get("Accept") != "application/json" {
		s.errorHandler(w, r, invalidMediaType, http.StatusUnsupportedMediaType)
		return true
	}
	return false
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

// isRunning returns a boolean representing whether the server is running or not.
func (s *Server) isRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
