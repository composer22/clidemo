package server

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	expectedInfoJSONResult = `{"version":"9.8.7","name":"Test Server","hostname":"localhost",` +
		`"UUID":"ABCDEFGHIJKLMNOPQRSTUVWXYZ","port":8080,"profPort":6060,"maxConnections":9999,` +
		`"maxWorkers":888,"debugEnabled":true}`
)

func TestInfoNew(t *testing.T) {
	info := InfoNew(func(i *Info) {
		i.Version = "9.8.7"
		i.Name = "Test Server"
		i.Hostname = "localhost"
		i.UUID = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		i.Port = 8080
		i.ProfPort = 6060
		i.MaxConn = 9999
		i.MaxWorkers = 888
		i.Debug = true
	})
	tp := reflect.TypeOf(info)

	if tp.Kind() != reflect.Ptr {
		t.Fatalf("Info not created as a pointer.\n")
	}

	tp = tp.Elem()
	if tp.Kind() != reflect.Struct {
		t.Fatalf("Info not created as a struct.\n")
	}
	if tp.Name() != "Info" {
		t.Fatalf("Info struct is not named correctly.\n")
	}
	if !(tp.NumField() > 0) {
		t.Fatalf("Info struct is empty.\n")
	}

}

func TestInfoString(t *testing.T) {
	t.Parallel()
	info := InfoNew(func(i *Info) {
		i.Version = "9.8.7"
		i.Name = "Test Server"
		i.Hostname = "localhost"
		i.UUID = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		i.Port = 8080
		i.ProfPort = 6060
		i.MaxConn = 9999
		i.MaxWorkers = 888
		i.Debug = true
	})
	if fmt.Sprint(info) != expectedInfoJSONResult {
		t.Errorf("Info not converted to json string.\n")
	}
}
