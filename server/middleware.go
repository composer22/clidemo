package server

import "net/http"

// middlewarePreprocess is used to perform filtering work on the request before the main controllers are called.
type middlewarePreprocess struct {
	serv    *Server
	handler http.Handler
}

// ServeHTTP implements the interface to accept requests so they can be filtered before handling by the server.
func (m *middlewarePreprocess) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.serv.incrementStats(r)
	m.serv.initResponseHeader(w)
	if m.serv.invalidHeader(w, r) {
		return
	}
	m.handler.ServeHTTP(w, r)
}
