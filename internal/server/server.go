package server

import (
	"net"
	"redis/pkg/transport"
	"sync"
)

type ConnectionHandler interface {
	Handle(*RedisConnection)
}

type RedisServer struct {
	transport transport.Transport

	handler ConnectionHandler

	connections         map[net.Addr]*RedisConnection
	connectionIdCounter int32
	sync.Mutex
}

func NewRedisServer(transport transport.Transport) *RedisServer {
	s := RedisServer{
		transport:   transport,
		connections: make(map[net.Addr]*RedisConnection),
	}

	return &s
}

func (s *RedisServer) Serve(handler ConnectionHandler) error {
	if err := s.transport.ListenAndAccept(); err != nil {
		return err
	}

	s.handler = handler

	s.startLoop()

	return nil
}

func (s *RedisServer) startLoop() {
	for conn := range s.transport.Consume() {
		go func(c net.Conn) {
			s.handleConnection(c)
		}(conn)
	}
}

// handleConnection runs in a separated goroutine
func (s *RedisServer) handleConnection(conn net.Conn) {
	s.Lock()

	s.connectionIdCounter++

	rconn := NewRedisConnection(conn, s.connectionIdCounter, 4096)
	s.connections[conn.RemoteAddr()] = rconn

	s.Unlock()

	s.handler.Handle(rconn)

	s.Lock()
	delete(s.connections, conn.RemoteAddr())
	s.Unlock()
}

func (s *RedisServer) Close() error {
	return s.transport.CLose()
}
