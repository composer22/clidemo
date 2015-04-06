package server

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

const (
	expectedStatsJSONResult = `{"startTime":"2006-01-02T13:24:56Z","requestCount":0,` +
		`"requestBytes":0,"connNumAvail":1234,"routeStats":{"route1":{"requesBytes":202,` +
		`"requestCounts":101},"route2":{"requesBytes":204,"requestCounts":103}}}`
)

func TestStatusNew(t *testing.T) {
	stats := StatusNew()
	tp := reflect.TypeOf(stats)

	if tp.Kind() != reflect.Ptr {
		t.Fatalf("Status not created as a pointer.")
	}

	tp = tp.Elem()
	if tp.Kind() != reflect.Struct {
		t.Fatalf("Status not created as a struct.")
	}
	if tp.Name() != "Status" {
		t.Fatalf("Status struct is not named correctly.")
	}
	if !(tp.NumField() > 0) {
		t.Fatalf("Status struct is empty.")
	}
}

func TestStatusIncrRequestStats(t *testing.T) {
	t.Parallel()
	stats := StatusNew()
	stats.IncrRequestStats(-1)
	if stats.RequestCount != 1 {
		t.Errorf("Status RequestCount not incremented correctly.")
	}
	if stats.RequestBytes != 0 {
		t.Errorf("Status RequestBytes should not have been incremented or decremented.")
	}

	stats.IncrRequestStats(101)
	stats.IncrRequestStats(99)
	if stats.RequestCount != 3 {
		t.Errorf("Status RequestCount not incremented correctly.")
	}
	if stats.RequestBytes != 200 {
		t.Errorf("Status RequestBytes should have been incremented.")
	}
}

func TestStatusIncrRouteStats(t *testing.T) {
	t.Parallel()
	stats := StatusNew()
	stats.IncrRouteStats("Route1", -1)

	rs, ok := stats.RouteStats["Route1"]
	if !ok {
		t.Errorf(`Status RouteStats["Route1"] entry not created correctly.`)
	}
	if len(rs) != 1 {
		t.Errorf(`Status RouteStats["Route1"] entry invalid size.`)
	}

	rc, ok := rs["requestCount"]
	if !ok {
		t.Errorf(`Status RouteStats["Route1"]["requestCount"] entry not created correctly.`)
	}
	if rc != 1 {
		t.Errorf(`Status RouteStats["Route1"]["requestCount"] should have been incremented.`)
	}

	rc, ok = rs["requestBytes"]
	if ok {
		t.Errorf(`Status RouteStats["Route1"]["requestBytes"] entry should not have been created.`)
	}

	stats = StatusNew()
	stats.IncrRouteStats("Route2", -1)
	stats.IncrRouteStats("Route2", 201)
	stats.IncrRouteStats("Route2", 98)
	if stats.RouteStats["Route2"]["requestCount"] != 3 {
		t.Errorf(`Status["Route2"]["requestCount"] not incremented correctly.`)
	}
	_, ok = stats.RouteStats["Route2"]["requestBytes"]
	if !ok {
		t.Errorf(`Status RouteStats["Route1"]["requestBytes"] entry should have been created.`)
	}
	if stats.RouteStats["Route2"]["requestBytes"] != 299 {
		t.Errorf(`Status["Route2"]["requestBytes"] not incremented correctly.`)
	}
}

func TestStatusString(t *testing.T) {
	t.Parallel()
	mockTime, _ := time.Parse(time.RFC1123Z, "Mon, 02 Jan 2006 13:24:56 -0000")
	stats := StatusNew(func(sts *Status) {
		sts.Start = mockTime
		sts.ConnNumAvail = 1234
		sts.RouteStats = map[string]map[string]int64{
			"route1": map[string]int64{"requestCounts": 101, "requesBytes": 202},
			"route2": map[string]int64{"requestCounts": 103, "requesBytes": 204},
		}
	})
	if fmt.Sprint(stats) != expectedStatsJSONResult {
		t.Errorf("Status not converted to json string.")
	}
}
