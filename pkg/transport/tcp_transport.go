package transport

import (
	"errors"
	"net"
	"redis/logger"
	"sync/atomic"
	"time"
)

var ErrServerClosed = errors.New("tcp: Server closed")

type tcpTransport struct {
	listenAddress string

	listener net.Listener

	closed atomic.Bool

	connChan chan net.Conn
}

func NewTcpTransport(listenAddress string) Transport {
	return &tcpTransport{
		listenAddress: listenAddress,

		connChan: make(chan net.Conn),
	}
}

func (t *tcpTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *tcpTransport) startAcceptLoop() {
	acceptErrorHandler := t.handleAcceptError()

	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			logger.Infof("tcp: server closed")
		}

		if err != nil {
			if err = acceptErrorHandler(err); err != nil {
				logger.Error(err)
				return
			}
		}

		go t.handleConnection(conn)
	}
}

func (t *tcpTransport) handleConnection(conn net.Conn) {
	t.connChan <- conn
}

func (t *tcpTransport) handleAcceptError() func(e error) error {
	var tempDelay time.Duration

	return func(err error) error {
		if t.shuttingDown() {
			return ErrServerClosed
		}

		var ne net.Error

		if errors.As(err, &ne) && ne.Temporary() {
			if tempDelay == 0 {
				tempDelay = 5 * time.Millisecond
			} else {
				tempDelay *= 2
			}
			if tempDelay > 1*time.Second {
				tempDelay = 1 * time.Second
			}

			logger.Errorf("redis: Accept error: %v; retrying in %v", err, tempDelay)

			time.Sleep(tempDelay)
			return nil
		}

		return err
	}
}

func (t *tcpTransport) shuttingDown() bool {
	return t.closed.Load()
}

func (t *tcpTransport) Consume() <-chan net.Conn {
	return t.connChan
}

func (t *tcpTransport) CLose() error {
	t.closed.Store(true)

	close(t.connChan)

	return t.listener.Close()
}
