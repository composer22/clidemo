package server

import (
	"encoding/json"
	"time"
)

// Status contains runtime statistics.
type Status struct {
	Start        time.Time                   `json:"startTime"`    // The start time of the server.
	RequestCount int64                       `json:"requestCount"` // How many requests came in to the server.
	RequestBytes int64                       `json:"requestBytes"` // Size of the requests in bytes.
	ConnNumAvail int                         `json:"connNumAvail"` // Number of live connections available.
	RouteStats   map[string]map[string]int64 `json:"routeStats"`   // How many requests/bytes came into each route.
}

// StatusNew is a factory function that returns a new instance of Status.
// options is an optional list of functions that initialize the structure
func StatusNew(options ...func(*Status)) *Status {
	st := &Status{
		Start:        time.Now(),
		ConnNumAvail: -1, // defaults to infinite.
		RouteStats:   make(map[string]map[string]int64),
	}
	for _, f := range options {
		f(st)
	}
	return st
}

// IncrRequestStats increments the stats totals for the server.
func (s *Status) IncrRequestStats(rb int64) {
	s.RequestCount++
	if rb > 0 {
		s.RequestBytes += rb
	}
}

// IncrRouteStats increments the stats totals for the route.
func (s *Status) IncrRouteStats(path string, rb int64) {
	if _, ok := s.RouteStats[path]; !ok {
		s.RouteStats[path] = make(map[string]int64)
	}

	s.RouteStats[path]["requestCount"]++
	if rb > 0 {
		s.RouteStats[path]["requestBytes"] += rb
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (s *Status) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}
