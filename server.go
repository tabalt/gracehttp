package gracehttp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	GRACEFUL_ENVIRON_KEY    = "IS_GRACEFUL"
	GRACEFUL_ENVIRON_STRING = GRACEFUL_ENVIRON_KEY + "=1"

	DEFAULT_READ_TIMEOUT  = 30 * time.Second
	DEFAULT_WRITE_TIMEOUT = DEFAULT_READ_TIMEOUT
)

// refer http.ListenAndServe
func ListenAndServe(addr string, handler http.Handler) error {
	return NewServer(addr, handler, DEFAULT_READ_TIMEOUT, DEFAULT_WRITE_TIMEOUT).ListenAndServe()
}

// refer http.ListenAndServeTLS
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	return NewServer(addr, handler, DEFAULT_READ_TIMEOUT, DEFAULT_WRITE_TIMEOUT).ListenAndServeTLS(certFile, keyFile)
}

// new server
func NewServer(addr string, handler http.Handler, readTimeout, writeTimeout time.Duration) *Server {

	// 获取环境变量
	isGraceful := false
	if os.Getenv(GRACEFUL_ENVIRON_KEY) != "" {
		isGraceful = true
	}

	// 实例化Server
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,

			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},

		isGraceful: isGraceful,
		signalChan: make(chan os.Signal),
	}
}

// 支持优雅重启的http服务
type Server struct {
	httpServer *http.Server
	listener   net.Listener

	isGraceful bool
	signalChan chan os.Signal
}

func (this *Server) ListenAndServe() error {
	addr := this.httpServer.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := this.getNetListener(addr)
	if err != nil {
		return err
	}

	this.listener = newListener(ln.(*net.TCPListener))

	return this.Serve()
}

func (this *Server) ListenAndServeTLS(certFile, keyFile string) error {
	addr := this.httpServer.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if this.httpServer.TLSConfig != nil {
		*config = *this.httpServer.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	ln, err := this.getNetListener(addr)
	if err != nil {
		return err
	}

	this.listener = tls.NewListener(newListener(ln.(*net.TCPListener)), config)
	return this.Serve()
}

func (this *Server) Serve() error {

	// 处理信号
	go this.handleSignals()

	// 处理HTTP请求
	err := this.httpServer.Serve(this.listener)

	// 跳出Serve处理代表 listener 已经close，等待所有已有的连接处理结束
	this.logf("listener closed, waiting for connection close...")
	this.listener.(*Listener).Wait()

	return err
}

func (this *Server) getNetListener(addr string) (ln net.Listener, err error) {

	if this.isGraceful {
		file := os.NewFile(3, "")
		ln, err = net.FileListener(file)
		if err != nil {
			err = fmt.Errorf("net.FileListener error: %v", err)
			return
		}
	} else {
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			err = fmt.Errorf("net.Listen error: %v", err)
			return
		}
	}
	return
}

func (this *Server) handleSignals() {
	var sig os.Signal

	signal.Notify(
		this.signalChan,
		syscall.SIGHUP,
		syscall.SIGTERM,
	)

	pid := os.Getpid()
	for {
		sig = <-this.signalChan

		switch sig {

		case syscall.SIGTERM:

			this.logf("pid %d received SIGTERM.", pid)
			this.logf("graceful shutting down http server...")
			this.shutdown()

		case syscall.SIGHUP:

			this.logf("pid %d received SIGHUP.", pid)
			this.logf("graceful restart http server...")

			err := this.fork()
			if err != nil {
				this.logf("fork error: %v.", err)
			} else {
				this.shutdown()
			}

		default:

		}
	}
}

// 优雅关闭
func (this *Server) shutdown() {

	// 通过设置超时使得进程不再接受新请求
	this.listener.(*Listener).SetDeadline(time.Now())

	// 关闭链接
	this.listener.(*Listener).Close()
}

// 启动子进程执行新程序
func (this *Server) fork() error {

	listenerFd, err := this.listener.(*Listener).GetFd()
	if err != nil {
		return fmt.Errorf("failed to get socket file descriptor: %v.", err)
	}

	path := os.Args[0]

	// 设置标识优雅重启的环境变量
	environList := []string{}
	for _, value := range os.Environ() {
		if value != GRACEFUL_ENVIRON_STRING {
			environList = append(environList, value)
		}
	}
	environList = append(environList, GRACEFUL_ENVIRON_STRING)

	execSpec := &syscall.ProcAttr{
		Env:   environList,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFd},
	}

	fork, err := syscall.ForkExec(path, os.Args, execSpec)
	if err != nil {
		return fmt.Errorf("failed to forkexec: %v.", err)
	}

	this.logf("fork exec to pid %d.", fork)

	return nil
}

func (this *Server) logf(format string, args ...interface{}) {

	if this.httpServer.ErrorLog != nil {
		this.httpServer.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
