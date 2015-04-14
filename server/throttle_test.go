package server

import "testing"

func TestThrottleConnClose(t *testing.T) {
	t.Parallel()
	// Covered by server test.
	t.SkipNow()
}

func TestThrottleConnDone(t *testing.T) {
	t.Parallel()
	t.Skip("Covered by server test.")
}

func TestThrottleListenerNew(t *testing.T) {
	t.Parallel()
	t.Skip("Covered by server test.")
}

func TestThrottleListenerAccept(t *testing.T) {
	t.Parallel()
	t.Skip("Covered by server test.")
}

func TestThrottleListenerGetConnectedCount(t *testing.T) {
	t.Parallel()
	t.Skip("Covered by server test.")
}
