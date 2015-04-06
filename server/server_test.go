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
	testSrvr = New(opts, func(s *Server) {})
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

func TestMethods(t *testing.T) {
	client := &http.Client{}

	req, _ := http.NewRequest("POST", "http://localhost:8080/v1.0/status", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body := strings.TrimSuffix(string(b), "\n")
	if body != InvalidMethod {
		t.Errorf("/status body should return method error.")
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("/status returned invalid method status code %d", resp.StatusCode)
	}

	req, _ = http.NewRequest("POST", "http://localhost:8080/v1.0/alive", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body = strings.TrimSuffix(string(b), "\n")
	if body != InvalidMethod {
		t.Errorf("/alive body should return method error.")
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("/alive returned invalid method status code %d", resp.StatusCode)
	}

	req, _ = http.NewRequest("GET", "http://localhost:8080/v1.0/parse",
		strings.NewReader(fmt.Sprintf(`{"text":"%s"}`, testParserText)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body = strings.TrimSuffix(string(b), "\n")
	if body != InvalidMethod {
		t.Errorf("/parse body should return method error.")
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("/parse returned invalid method status code %d", resp.StatusCode)
	}
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
	if resp.StatusCode != http.StatusOK {
		t.Errorf("/status returned invalid status code %d", resp.StatusCode)
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
	if resp.StatusCode != http.StatusOK {
		t.Errorf("/alive returned invalid status code %d", resp.StatusCode)
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
		t.Errorf("Body should not be empty.")
	}
	if body != testParserResultJSON {
		t.Errorf("/parse returned invalid results: %s", body)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Invalid /parse status code %d", resp.StatusCode)
	}
}

func TestParseHandler(t *testing.T) {
	client := &http.Client{}

	req, _ := http.NewRequest("POST", "http://localhost:8080/v1.0/parse",
		strings.NewReader(fmt.Sprintf(`"text":"%s"`, testParserText)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body := strings.TrimSuffix(string(b), "\n")
	if body != InvalidJSONText {
		t.Errorf("JSON body should have been found invalid.")
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("/parse status code incorrect for bad JSON: %d", resp.StatusCode)
	}

	req, _ = http.NewRequest("POST", "http://localhost:8080/v1.0/parse",
		strings.NewReader(`{"monkey":"monkeytext"}`))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 3A3E6C4C51F12DF2415682CCF9D18")
	resp, _ = client.Do(req)
	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body = strings.TrimSuffix(string(b), "\n")
	if body != InvalidJSONAttribute {
		t.Errorf("JSON attr should have been found invalid.")
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("/parse status code incorrect for bad JSON attr: %d", resp.StatusCode)
	}
}

func TestServerTakeDown(t *testing.T) {

	time.Sleep(2 * time.Second) // Coverage of timeout in Throttle.
	testSrvr.Shutdown()
	testSrvr.Shutdown() // Coverage of isRunning test in Shutdown().
	if testSrvr.isRunning() {
		t.Errorf("Server should have shut down.")
	}
	testSrvr = nil
}
