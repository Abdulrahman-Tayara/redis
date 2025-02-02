package commands

import (
	"errors"
	"redis/internal/server"
	"redis/pkg/iox"
)

func (s *Server) HandleGet() server.CommandHandlerFunc {
	return func(ctx *server.Context, w iox.AnyWriter) {
		args := ctx.Args()

		if len(args) < 1 {
			_, _ = w.WriteAny(errors.New("ERR wrong number of arguments for 'get' command"))
			return
		}

		key := args[0]

		if v, err := s.store.Get(key.(string)); err != nil {
			_, _ = w.WriteAny(err)
			return
		} else {
			_, _ = w.WriteAny(v)
		}
	}
}
