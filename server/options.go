package server

// Options represents parameters that are passed to the application to be used in constructing the run and
// the server (if server mode is indicated).
type Options struct {
	Name     string `json:"name"`           // The name of the server.
	Port     int    `json:"port"`           // The default port of the server.
	MaxConn  int    `json:"maxConnections"` // The maximum concurrent connections accepted by the server.
	MaxProcs int    `json:"maxProcs"`       // The maximum number of processor cores available for the application/server.
	Debug    bool   `json:"debugEnabled"`   // Is debugging enabled in the application or server.
}
