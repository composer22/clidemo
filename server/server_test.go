package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"testing"
	"time"
)

const (
	testParserText       = "Now is the 'Winter' of our discontent. And then the other dude as well."
	testParserResultJSON = `{"result":{"words":{"and":{"counter":1,"sentenceUse":[1]},"as":` +
		`{"counter":1,"sentenceUse":[1]},"discontent":{"counter":1,"sentenceUse":[0]},` +
		`"dude":{"counter":1,"sentenceUse":[1]},"is":{"counter":1,"sentenceUse":[0]},` +
		`"now":{"counter":1,"sentenceUse":[0]},"of":{"counter":1,"sentenceUse":[0]},` +
		`"other":{"counter":1,"sentenceUse":[1]},"our":{"counter":1,"sentenceUse":[0]},` +
		`"the":{"counter":2,"sentenceUse":[0,1]},"then":{"counter":1,"sentenceUse":[1]},` +
		`"well":{"counter":1,"sentenceUse":[1]},"winter":{"counter":1,"sentenceUse":[0]}}}}`
)

var (
	testSrvr *Server
)

func TestServerStartup(t *testing.T) {
	opts := &Options{
		Name:       "Test Server",
		Hostname:   "localhost",
		Port:       8080,
		ProfPort:   6060,
		MaxConn:    1000,
		MaxWorkers: 1000,
		MaxProcs:   1,
		Debug:      true,
	}

	runtime.GOMAXPROCS(1)
	testSrvr = New(opts)
	go func() { testSrvr.Start() }()
}

func TestValidHeaders(t *testing.T) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://localhost:8080/v1.0/alive", nil)
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body := strings.TrimSuffix(string(b), "\n")
	if body != InvalidMediaType {
		t.Errorf("Missing 'Accept' header returned invalid body: %s", body)
	}
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("Missing 'Accept' header returned invalid  status code %d", resp.StatusCode)
	}

	req.Header.Add("Accept", "text/html")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body = strings.TrimSuffix(string(b), "\n")
	if body != InvalidMediaType {
		t.Errorf("Invalid 'Accept' header returned invalid body: %s", body)
	}
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("Invalid 'Accept' header returned invalid  status code %d", resp.StatusCode)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Del("Content-Type")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body = strings.TrimSuffix(string(b), "\n")
	if body != InvalidMediaType {
		t.Errorf("Missing 'Content-Type' header returned invalid body: %s", body)
	}
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("Missing 'Content-Type' header returned invalid  status code %d", resp.StatusCode)
	}

	req.Header.Add("Content-Type", "text/html")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body = strings.TrimSuffix(string(b), "\n")
	if body != InvalidMediaType {
		t.Errorf("Invalid 'Content-Type' header returned invalid body: %s", body)
	}
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("Invalid 'Content-Type' header returned invalid  status code %d", resp.StatusCode)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Del("Authorization")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body = strings.TrimSuffix(string(b), "\n")
	if body != InvalidAuthorization {
		t.Errorf("Missing 'Authorization' header returned invalid body: %s", body)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Missing 'Authorization' header returned invalid  status code %d", resp.StatusCode)
	}

	req.Header.Add("Authorization", "Bearer XXX3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body = strings.TrimSuffix(string(b), "\n")
	if body != InvalidAuthorization {
		t.Errorf("Invalid 'Authorization' header returned invalid body: %s", body)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Invalid 'Authorization' header returned invalid  status code %d", resp.StatusCode)
	}
	req.Header.Set("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
}

func TestRoutes(t *testing.T) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://localhost:8080/v1.0/status", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body := string(b)
	if body == "" {
		t.Errorf("/status body should not be empty.")
	}
	if resp.StatusCode != 200 {
		t.Errorf("/status returned invalid  status code %d", resp.StatusCode)
	}

	req, _ = http.NewRequest("GET", "http://localhost:8080/v1.0/alive", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body = string(b)
	if body != "" {
		t.Errorf("/alive body should be empty.")
	}
	if resp.StatusCode != 200 {
		t.Errorf("/alive returned invalid  status code %d", resp.StatusCode)
	}

	req, _ = http.NewRequest("POST", "http://localhost:8080/v1.0/parse",
		strings.NewReader(fmt.Sprintf(`{"text":"%s"}`, testParserText)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body = string(b)
	if body == "" {
		t.Errorf("Body should not be empty\n")
	}
	if body != testParserResultJSON {
		t.Errorf("/parse returned invalid results: %s", body)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Invalid /parse status code %d", resp.StatusCode)
	}
}

func TestServerTakeDown(t *testing.T) {
	testSrvr.Shutdown()
	testSrvr = nil

	var maxTimeout time.Duration
	if TCPReadTimeout > TCPWriteTimeout {
		maxTimeout = TCPReadTimeout
	} else {
		maxTimeout = TCPWriteTimeout
	}
	time.Sleep(maxTimeout + (1 * time.Second))
}
