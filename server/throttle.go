package server

import (
	"net"
	"sync"
)

// ThrottledConn is a decorator over net.conn that allows us to throttle connections via the listener.
type ThrottledConn struct {
	*net.TCPConn
	acceptCh  chan bool
	closeOnce sync.Once
}

// Close overloads the class function of the connection so that the listener throttle can be serviced.
func (c *ThrottledConn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		defer c.Done()
		err = c.TCPConn.Close()
	})
	return err
}

// Done puts back a token so it can be serviced again by the throttle listener.
func (c *ThrottledConn) Done() {
	c.acceptCh <- true // Put back token.
}

// ThrottledListener is a wrapper on a listener that limits connections.
type ThrottledListener struct {
	*net.TCPListener
	acceptCh chan bool // Queue for service tokens
	maxConns int
}

// ThrottledListenerNew is a factory function that returns an instatiated ThrottledListener.
func ThrottledListenerNew(addr string, maxConnAllowed int) (*ThrottledListener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Initialize accept tokens.
	acceptCh := make(chan bool, maxConnAllowed)
	for i := 0; i < maxConnAllowed; i++ {
		acceptCh <- true
	}

	return &ThrottledListener{
		TCPListener: ln.(*net.TCPListener),
		acceptCh:    acceptCh,
		maxConns:    maxConnAllowed,
	}, nil
}

// Accept overrides the accept function of the listener so that waits can occur on tokens in the queue.
func (t *ThrottledListener) Accept() (net.Conn, error) {

	// Wait to grab a token and service a connection.
	<-t.acceptCh
	conn, err := t.TCPListener.AcceptTCP()
	if err != nil {
		t.acceptCh <- true // err so put token back
		return nil, err
	}

	// Set connection to stay alive and return it.
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(TCPConnectionTimeout)
	return &ThrottledConn{
		TCPConn:  conn,
		acceptCh: t.acceptCh,
	}, nil
}

// GetConnectedCount returns the total number of connections at present being serviced.
func (t *ThrottledListener) GetConnectedCount() int {
	return t.maxConns - len(t.acceptCh)
}
