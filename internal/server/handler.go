package server

import (
	"context"
	"errors"
	"net"
	"redis/logger"
	"redis/pkg/iox"
	"redis/pkg/resp"
	"strings"
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

type CommandHandlerFunc func(ctx *Context, w iox.AnyWriter)

func (c CommandHandlerFunc) Handle(ctx *Context, w iox.AnyWriter) {
	c(ctx, w)
}

type CommandHandler interface {
	Handle(ctx *Context, w iox.AnyWriter)
}

type ConnectionServe struct {
	commandHandlers map[string]CommandHandler
}

func NewConnectionServe() *ConnectionServe {
	return &ConnectionServe{
		commandHandlers: make(map[string]CommandHandler),
	}
}

func (h *ConnectionServe) Command(command string, handler CommandHandler) {
	h.commandHandlers[h.normalizeCommand(command)] = handler
}

func (h *ConnectionServe) Handle(conn *RedisConnection) {
	defer func() {
		h.close(conn)
	}()

	for {
		command, args, err := conn.ReadCommand()
		if err != nil {
			// the resp.ErrReaderRead belongs to the connection.Read() error
			if errors.Is(resp.ErrReaderRead, err) || errors.Is(resp.ErrReaderRead, errors.Unwrap(err)) {
				break
			}

			h.writeError(err, conn)
			continue
		}

		command = h.normalizeCommand(command)

		commandHandler, ok := h.commandHandlers[command]
		if !ok {
			logger.Errorf("command %s is not found, args: %v", command, args)
			h.writeError(ErrCommandNotFound, conn)
			continue
		}

		ctx := newContext(context.TODO(), conn, command, args)

		commandHandler.Handle(ctx, conn)
	}
}

func (h *ConnectionServe) normalizeCommand(cmd string) string {
	return strings.ToLower(cmd)
}

func (h *ConnectionServe) writeError(err error, w iox.AnyWriter) {
	logger.Error(err)
	if _, err = w.WriteAny(err); err != nil {
		logger.Errorf("conn write err: %v", err.Error())
	}
}

func (h *ConnectionServe) close(conn *RedisConnection) {
	logger.Infof("closing %s connection", conn.RemoteAddr().String())

	if err := conn.Close(); err != nil {
		logger.Errorf("conn close err: %v", err.Error())
	}
}
