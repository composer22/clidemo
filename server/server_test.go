package server

import (
	"io/ioutil"
	"net/http"
	"runtime"
	"testing"
)

func TestRoutes(t *testing.T) {
	opts := &Options{
		Name:     "Test Server",
		Hostname: "localhost",
		Port:     8080,
		ProfPort: 6060,
		MaxConn:  1000,
		MaxProcs: 1000,
		Debug:    true,
	}
	runtime.GOMAXPROCS(1)
	srvr := New(opts)
	go func() { srvr.Start() }()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/v1.0/status", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)
	body := string(b)
	if body == "" {
		t.Errorf("Invalid Body\n")
	}
	if resp.StatusCode != 200 {
		t.Errorf("Invalid /status status code %d\n", resp.StatusCode)
	}

	req, _ = http.NewRequest("GET", "http://localhost:8080/v1.0/alive", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	body = string(b)
	if body != "" {
		t.Errorf("Body should be empty\n")
	}
	if resp.StatusCode != 200 {
		t.Errorf("Invalid /alive status code %d\n", resp.StatusCode)
	}

	srvr.Shutdown()
}
