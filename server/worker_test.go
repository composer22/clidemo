package server

import (
	"sync"
	"testing"
)

const (
	workerParseTestText      = "This is a test. This is another test."
	expectedWorkerJSONResult = `{"words":{"a":{"counter":1,"sentenceUse":[0]},"another":` +
		`{"counter":1,"sentenceUse":[1]},"is":{"counter":2,"sentenceUse":[0,1]},` +
		`"test":{"counter":2,"sentenceUse":[0,1]},"this":{"counter":2,"sentenceUse":[0,1]}}}`
)

func TestParseWorker(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup
	jobq := make(chan *parseJob)
	wg.Add(1)
	go parseWorker(jobq, &wg)
	job := parseJob{
		Source: workerParseTestText,
		DoneCh: make(chan bool),
	}
	jobq <- &job
	<-job.DoneCh
	if job.Result != expectedWorkerJSONResult {
		t.Errorf("Worker expected parse doesn't match result.")
	}
	close(jobq)
	wg.Wait()
}
