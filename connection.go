package gracehttp

import (
	"net"
)

type Connection struct {
	net.Conn
	listener *Listener
}

func (this *Connection) Close() error {
	this.listener.waitGroup.Done()
	return this.Conn.Close()
}
