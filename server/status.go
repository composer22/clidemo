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
	CurrentConns int                         `json:"currentConns"` // Current number of live connections.
	RouteStats   map[string]map[string]int64 `json:"routeStats"`   // How many requests/bytes came into each route.
}

// IncrRequestStats increments the stats totals for the server.
func (s *Status) IncrRequestStats(reqBytes int64) {
	s.RequestCount++
	if reqBytes > 0 {
		s.RequestBytes += reqBytes
	}
}

// IncrRouteStats increments the stats totals for the route.
func (s *Status) IncrRouteStats(path string, reqBytes int64) {
	if _, ok := s.RouteStats[path]; !ok {
		s.RouteStats[path] = make(map[string]int64)
	}

	s.RouteStats[path]["requestCount"]++
	if reqBytes > 0 {
		s.RouteStats[path]["requestBytes"] += reqBytes
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a string to fmt.Print() etc.
func (s *Status) String() string {
	result, _ := json.Marshal(s)
	return string(result)
}
