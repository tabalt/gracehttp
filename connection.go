package gracehttp

import (
	"net"
)

type Connection struct {
	net.Conn
	listener *Listener

	closed bool
}

func (this *Connection) Close() error {

	if !this.closed {
		this.closed = true
		this.listener.waitGroup.Done()
	}

	return this.Conn.Close()
}
