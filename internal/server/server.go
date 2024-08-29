package server

import (
	"context"
	"errors"
	"net"
	"redis/internal/transport"
	"redis/logger"
	"redis/pkg/resp"
)

type RedisServer struct {
	transport transport.Transport

	handlers map[string]CommandHandler
}

func NewRedisServer(transport transport.Transport) *RedisServer {
	s := RedisServer{
		transport: transport,
		handlers:  make(map[string]CommandHandler),
	}

	return &s
}

func (s *RedisServer) Handle(command string, handler CommandHandler) {
	s.handlers[command] = handler
}

func (s *RedisServer) Serve() error {
	if err := s.transport.ListenAndAccept(); err != nil {
		return err
	}

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
	defer func() {
		logger.Infof("closing %s connection", conn.RemoteAddr().String())

		if err := conn.Close(); err != nil {
			logger.Errorf("conn close err: %v", err.Error())
		}
	}()

	w := resp.NewRespWriter(conn)

	writeErrFunc := func(err error) {
		if _, err = w.WriteAny(err); err != nil {
			logger.Errorf("conn write err: %v", err.Error())
		}
	}

	for {
		buf := make([]byte, 4096)

		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		command, args, err := ReadRespCommand(buf[:n])
		if err != nil {
			logger.Error(err)
			writeErrFunc(errors.New("INVALID_RESP_CONTENT"))
			continue
		}

		commandHandler, ok := s.handlers[command]
		if !ok {
			writeErrFunc(errors.New("INVALID_COMMAND"))
			continue
		}

		ctx := newContext(context.TODO(), conn, command, args)

		commandHandler.Handle(ctx, w)
	}

}

func (s *RedisServer) Close() error {
	return s.transport.CLose()
}
