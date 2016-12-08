package gracehttp

import (
	"net"
	"sync"
	"time"
)

func newListener(ln *net.TCPListener) net.Listener {
	return &Listener{
		TCPListener: ln,
		waitGroup:   &sync.WaitGroup{},
	}
}

type Listener struct {
	*net.TCPListener

	waitGroup *sync.WaitGroup
}

func (this *Listener) Accept() (net.Conn, error) {

	tc, err := this.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(time.Minute)

	this.waitGroup.Add(1)

	conn := &Connection{
		Conn:     tc,
		listener: this,
	}
	return conn, nil
}

func (this *Listener) Wait() {
	this.waitGroup.Wait()
}

func (this *Listener) GetFd() (uintptr, error) {
	file, err := this.TCPListener.File()
	if err != nil {
		return 0, err
	}
	return file.Fd(), nil
}
