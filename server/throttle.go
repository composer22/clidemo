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
	wg        sync.WaitGroup
	acceptCh  chan bool
	closeOnce sync.Once
}

// Close overloads the class function of the connection so that the listener throttle can be serviced.
func (c *ThrottledConn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		defer c.wg.Done()
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
	wg       sync.WaitGroup // For waiting on connections to close
	acceptCh chan bool      // Queue for service tokens.
	stopCh   chan bool      // Shutdown server requested.
	maxConns int
}

// ThrottledListenerNew is a factory function that returns an instatiated ThrottledListener.
func ThrottledListenerNew(addr string, maxConnAllowed int) (*ThrottledListener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Initialize accept tokens.
	var acceptCh chan bool
	if maxConnAllowed > 0 {
		acceptCh = make(chan bool, maxConnAllowed)
		for i := 0; i < maxConnAllowed; i++ {
			acceptCh <- true
		}
	}
	return &ThrottledListener{
		TCPListener: ln.(*net.TCPListener),
		acceptCh:    acceptCh,
		stopCh:      make(chan bool),
		maxConns:    maxConnAllowed,
	}, nil
}

// Accept overrides the accept function of the listener so that waits can occur on tokens in the queue.
func (t *ThrottledListener) Accept() (net.Conn, error) {
	for {
		// Wait to grab a token and service a connection.
		if t.acceptCh != nil {
			<-t.acceptCh
		}

		// Look for a request
		t.SetDeadline(time.Now().Add(time.Second))
		conn, err := t.TCPListener.AcceptTCP()
		if err != nil {
			if t.acceptCh != nil {
				t.acceptCh <- true // err so put token back
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() && netErr.Temporary() {
				continue // Its a timeout, so try again...
			}
			return nil, err
		}

		// Check for shutdown signal
		select {
		case <-t.stopCh:
			t.wg.Wait()
			t.Close()
			return nil, StoppedError
		default: // continue
		}

		// Set connection to stay alive and return it.
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(TCPConnectionTimeout)
		t.wg.Add(1)
		return &ThrottledConn{
			TCPConn:  conn,
			wg:       t.wg,
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
