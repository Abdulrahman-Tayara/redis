package commands

import (
	"redis/internal/server"
	"redis/pkg/iox"
)

func (s *Server) HandleClientCommands() server.CommandHandlerFunc {
	return func(ctx *server.Context, w iox.AnyWriter) {
		_, _ = w.WriteAny("OK")
	}
}
