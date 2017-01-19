package gracehttp

import (
	"net"
	"sync"
	"time"
)

type Listener struct {
	*net.TCPListener

	wg *sync.WaitGroup
}

func NewListener(tl *net.TCPListener) net.Listener {
	return &Listener{
		TCPListener: tl,

		wg: &sync.WaitGroup{},
	}
}

func (l *Listener) Fd() (uintptr, error) {
	file, err := l.TCPListener.File()
	if err != nil {
		return 0, err
	}
	return file.Fd(), nil
}

func (l *Listener) Accept() (net.Conn, error) {

	tc, err := l.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(time.Minute)

	l.wg.Add(1)

	conn := &Connection{
		Conn:     tc,
		listener: l,
	}
	return conn, nil
}

func (l *Listener) Wait() {
	l.wg.Wait()
}
