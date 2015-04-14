package server

import (
	"errors"
	"net"
	"sync"
	"time"
)

var (
	StoppedError = errors.New("Server stop requested.")
)

// ThrottledConn is a wrapper over net.conn that allows us to throttle connections via the listener.
type ThrottledConn struct {
	*net.TCPConn
	acceptCh  chan bool
	closeOnce sync.Once
}

// Close overloads the type function of the connection so that the listener throttle can be serviced.
// TODO If file streaming is needed in the future, also add CloseRead() and CloseWrite() coverage.
func (c *ThrottledConn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		c.Done()
		err = c.TCPConn.Close()
	})
	return err
}

// Done puts back a token so it can be serviced again by the throttle listener.
func (c *ThrottledConn) Done() {
	if c.acceptCh != nil {
		c.acceptCh <- true // Put back token.
	}
}

// ThrottledListener is a wrapper on a listener that limits connections.
type ThrottledListener struct {
	*net.TCPListener
	acceptCh chan bool // Queue for service tokens.
	stopCh   chan bool // Shutdown server requested.
	maxConns int
}

// ThrottledListenerNew is a factory function that returns an instatiated ThrottledListener.
func ThrottledListenerNew(addr string, mxConn int) (*ThrottledListener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Initialize accept tokens.
	var acceptCh chan bool
	if mxConn > 0 {
		acceptCh = make(chan bool, mxConn)
		for i := 0; i < mxConn; i++ {
			acceptCh <- true
		}
	}
	return &ThrottledListener{
		TCPListener: ln.(*net.TCPListener),
		acceptCh:    acceptCh,
		stopCh:      make(chan bool),
		maxConns:    mxConn,
	}, nil
}

// Accept overrides the accept function of the listener so that waits can occur on
// tokens in the queue.
func (t *ThrottledListener) Accept() (net.Conn, error) {
	for {
		// Wait to grab a token if we are in restricted mode.
		if t.acceptCh != nil {
			<-t.acceptCh
		}

		// Look for a request for one second.
		t.SetDeadline(time.Now().Add(time.Second))
		conn, err := t.AcceptTCP()

		// Check for shutdown signal
		select {
		case <-t.stopCh:
			t.Close()
			return nil, StoppedError
		default: // continue
		}

		if err != nil {
			// Return token if we are in restricted mode.
			if t.acceptCh != nil {
				t.acceptCh <- true
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
			return nil, err
		}

		// Set connection to stay alive n-time and return it.
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(TCPKeepAliveTimeout)
		return &ThrottledConn{
			TCPConn:  conn,
			acceptCh: t.acceptCh,
		}, nil
	}
}

// Stops the listener
func (t *ThrottledListener) Stop() {
	close(t.stopCh)
}

// GetConnNumAvail returns the total number of connections available.
func (t *ThrottledListener) GetConnNumAvail() int {
	if t.acceptCh != nil {
		return len(t.acceptCh)
	}
	return -1
}
