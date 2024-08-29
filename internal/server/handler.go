package server

import (
	"context"
	"net"
	"redis/pkg/iox"
)

type Context struct {
	context.Context

	conn net.Conn

	command string
	args    []any
}

func newContext(ctx context.Context, conn net.Conn, command string, args []any) *Context {
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

type CommandHandler interface {
	Handle(ctx *Context, w iox.AnyWriter)
}
