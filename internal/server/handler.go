package server

import (
	"context"
	"errors"
	"io"
	"net"
	"redis/logger"
	"redis/pkg/iox"
)

var (
	ErrCommandNotFound = errors.New("COMMAND_NOT_FOUND")
)

type Context struct {
	context.Context

	conn *RedisConnection

	command string
	args    []any
}

func newContext(ctx context.Context, conn *RedisConnection, command string, args []any) *Context {
	return &Context{
		Context: ctx,
		conn:    conn,
		command: command,
		args:    args,
	}
}

func (c *Context) Args() []any {
	return c.args
}

func (c *Context) Command() string {
	return c.command
}

func (c *Context) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Context) ConnectionId() int32 {
	return c.conn.Id
}

func (c *Context) Set(key, value string) {
	c.conn.SetInfo(key, value)
}

func (c *Context) Get(key string) (string, bool) {
	return c.conn.GetInfo(key)
}

type CommandHandlerFunc func(ctx *Context, w iox.AnyWriter)

func (c CommandHandlerFunc) Handle(ctx *Context, w iox.AnyWriter) {
	c(ctx, w)
}

type CommandHandler interface {
	Handle(ctx *Context, w iox.AnyWriter)
}

type ConnectionServerOptions struct {
	ErrHandler    func(err error, w iox.AnyWriter)
	CommandMapper func(string) string
}

type ConnectionServe struct {
	commandHandlers map[string]CommandHandler

	errHandler func(err error, w iox.AnyWriter)

	commandMapper func(string) string
}

func NewConnectionServe(opts *ConnectionServerOptions) *ConnectionServe {
	server := &ConnectionServe{
		commandHandlers: make(map[string]CommandHandler),
		errHandler:      writeError,
		commandMapper: func(s string) string {
			return s
		},
	}

	if opts != nil {
		if opts.ErrHandler != nil {
			server.errHandler = opts.ErrHandler
		}
		if opts.CommandMapper != nil {
			server.commandMapper = opts.CommandMapper
		}
	}

	return server
}

func (h *ConnectionServe) Command(command string, handler CommandHandler) {
	h.commandHandlers[h.commandMapper(command)] = handler
}

func (h *ConnectionServe) Handle(conn *RedisConnection) {
	defer func() {
		h.close(conn)
	}()

	for {
		command, args, err := conn.ReadCommand()
		if err != nil {
			if err == io.EOF {
				break
			}

			// the resp.ErrReaderRead belongs to the connection.Read() error
			var netErr net.Error
			if errors.As(err, &netErr) {
				logger.Errorf("conn read err: %v", err.Error())
				break
			}

			h.errHandler(err, conn)
			continue
		}
		if command == "" {
			continue
		}

		command = h.commandMapper(command)

		logger.Infof("executing command: %s, args: %v", command, args)

		commandHandler, ok := h.commandHandlers[command]
		if !ok {
			logger.Errorf("command %s is not found, args: %v", command, args)
			h.errHandler(ErrCommandNotFound, conn)
			continue
		}

		ctx := newContext(context.TODO(), conn, command, args)

		commandHandler.Handle(ctx, conn)
	}
}

func (h *ConnectionServe) close(conn *RedisConnection) {
	logger.Infof("closing %s connection", conn.RemoteAddr().String())

	if err := conn.Close(); err != nil {
		logger.Errorf("conn close err: %v", err.Error())
	}
}

func writeError(err error, w iox.AnyWriter) {
	logger.Error(err)
	if _, err = w.WriteAny(err); err != nil {
		logger.Errorf("conn write err: %v", err.Error())
	}
}
