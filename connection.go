package gracehttp

import (
	"net"
)

type Connection struct {
	net.Conn
	listener *Listener

	closed bool
}

func (conn *Connection) Close() error {

	if !conn.closed {
		conn.closed = true
		conn.listener.wg.Done()
	}

	return conn.Conn.Close()
}
