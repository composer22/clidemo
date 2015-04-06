package server

import "encoding/json"

// Info provides basic information about the running server.
type Info struct {
	Version    string `json:"version"`        // Version of the server.
	Name       string `json:"name"`           // The name of the server.
	Hostname   string `json:"hostname"`       // The hostname of the server.
	UUID       string `json:"UUID"`           // Unique ID of the server.
	Port       int    `json:"port"`           // Port the server is listening on.
	ProfPort   int    `json:"profPort"`       // Profiler port the server is listening on.
	MaxConn    int    `json:"maxConnections"` // The maximum concurrent connections accepted.
	MaxWorkers int    `json:"maxWorkers"`     // The maximum numer of workers allowed to run.
	Debug      bool   `json:"debugEnabled"`   // Is debugging enabled on the server.
}

// InfoNew is a factory function that returns a new instance of Info.
// options is an optional list of functions that initialize the structure
func InfoNew(options ...func(*Info)) *Info {
	inf := &Info{
		Version: version,
		UUID:    createV4UUID(),
	}
	for _, option := range options {
		option(inf)
	}
	return inf
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (i *Info) String() string {
	result, _ := json.Marshal(i)
	return string(result)
}
