package commands

import (
	"fmt"
	"redis/internal/server"
	"redis/pkg/iox"
	"strconv"
)

func (s *Server) HandleHello() server.CommandHandlerFunc {
	return func(ctx *server.Context, w iox.AnyWriter) {
		items := map[string]any{
			"# server":  quotation("redis"),
			"# version": quotation(s.cfg.Version),
			"# proto":   strconv.Itoa(s.cfg.ProtoVersion),
			"# id":      strconv.Itoa(int(ctx.ConnectionId())),
			"# mode":    quotation(s.cfg.Mode),
			"# role":    quotation("master"),
			"# modules": "(empty array)",
		}

		_, _ = w.WriteAny(items)
	}
}

func quotation(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}
