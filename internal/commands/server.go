package commands

import (
	"redis/internal/configs"
	"redis/internal/server"
	"redis/internal/store"
	"strings"
)

type Server struct {
	cfg *configs.Configs

	store store.Store
}

func NewServer(cfg *configs.Configs, store store.Store) *Server {
	return &Server{
		cfg:   cfg,
		store: store,
	}
}

func (s *Server) Handlers() server.ConnectionHandler {
	handler := server.NewConnectionServe(&server.ConnectionServerOptions{
		CommandMapper: normalizeCommand,
	})

	handler.Command("hello", s.HandleHello())
	handler.Command("set", s.HandleSet())
	handler.Command("get", s.HandleGet())
	handler.Command("client", s.HandleClientCommands())

	return handler
}

func normalizeCommand(cmd string) string {
	return strings.ToLower(cmd)
}
