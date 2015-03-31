package server

import "net/http"

// Middleware is used to perform filtering work on the request before the main controllers are called.
type Middleware struct {
	serv    *Server
	handler http.Handler
}

// ServeHTTP implements the interface to accept requests so they can be filtered before handling by the server.
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.serv.LogRequest(r)
	m.serv.incrementStats(r)
	m.serv.initResponseHeader(w)
	if m.serv.invalidHeader(w, r) {
		return
	}
	m.handler.ServeHTTP(w, r)
}
