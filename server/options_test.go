package server

import (
	"fmt"
	"testing"
)

const (
	expectedOptionsJSONResult = `{"name":"Test Options","hostname":"localhost","port":8080,` +
		`"profPort":6060,"maxConnections":1001,"maxWorkers":999,"maxProcs":888,` +
		`"debugEnabled":true}`
)

func TestOptionsString(t *testing.T) {
	options := &Options{
		Name:       "Test Options",
		Hostname:   "localhost",
		Port:       8080,
		ProfPort:   6060,
		MaxConn:    1001,
		MaxWorkers: 999,
		MaxProcs:   888,
		Debug:      true,
	}
	if fmt.Sprint(options) != expectedOptionsJSONResult {
		t.Errorf("Options not converted to json string.\n")
	}
}
