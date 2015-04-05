package test

import (
	"io/ioutil"
	"net/http"
	"runtime"
	"testing"

	"github.com/composer22/clidemo/server"
)

type Options struct {
	Name       string `json:"name"`           // The name of the server.
	Hostname   string `json:"hostname"`       // The hostname of the server.
	Port       int    `json:"port"`           // The default port of the server.
	ProfPort   int    `json:"profPort"`       // The profiler port of the server.
	MaxConn    int    `json:"maxConnections"` // The maximum concurrent connections accepted.
	MaxWorkers int    `json:"maxWorkers"`     // The maximum numer of workers allowed to run.
	MaxProcs   int    `json:"maxProcs"`       // The maximum number of processor cores available.
	Debug      bool   `json:"debugEnabled"`   // Is debugging enabled in the application or server.
}

func TestRoutes(t *testing.T) {
	t.Parallel()
	opts := &server.Options{
		Name:     "Test Server",
		Hostname: "localhost",
		Port:     8080,
		ProfPort: 6060,
		MaxConn:  1000,
		MaxProcs: 1000,
		Debug:    true,
	}
	runtime.GOMAXPROCS(1)
	srvr := server.New(opts)

	go func() { srvr.Start() }()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/v1.0/status", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)
	t.Logf(`"code":"%d","body":"%s"`, resp.StatusCode, string(b))

	req, _ = http.NewRequest("GET", "http://localhost:8080/v1.0/alive", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	t.Logf(`"code":"%d","body":"%s"`, resp.StatusCode, string(b))

	srvr.Shutdown()
}
