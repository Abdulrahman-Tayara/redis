package commands

import (
	"redis/internal/configs"
	"redis/internal/server"
	"redis/internal/store"
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
	handler := server.NewConnectionServe()

	handler.Command("hello", s.HandleHello())
	handler.Command("set", s.HandleSet())

	return handler
}
