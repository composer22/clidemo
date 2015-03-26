package server

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/composer22/clidemo/parser"
)

// parseJob is a transport packet that represents text that needs parsing by a worker.
type parseJob struct {
	Source     string    `json:"source"` // Source text to be parsed.
	DoneCh     chan bool `json:"-"`      // Channel to notify when done parsing.
	ResultJSON string    `json:"result"` // Result JSON string of parse results
}

// parseWorker is used as a go routine wrapper to handle parsing jobs for the server.
func parseWorker(jobq chan *parseJob, wg *sync.WaitGroup) {
	defer wg.Done()
	p := parser.New()
	for {
		job, ok := <-jobq
		if !ok {
			break
		}
		p.Execute(bytes.NewBufferString(job.Source))
		job.ResultJSON = fmt.Sprint(p)
		job.DoneCh <- true
		p.Reset()
	}
}
