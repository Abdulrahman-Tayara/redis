package main

import (
	"redis/internal/commands"
	"redis/internal/configs"
	"redis/internal/server"
	store2 "redis/internal/store"
	"redis/pkg/transport"
)

func main() {

	store := store2.NewInMemoryStore()

	commandsServer := commands.NewServer(&configs.Configs{
		Version:      "6.0.3",
		ProtoVersion: 3,
		Mode:         "standalone",
		Modules:      []string{},
	}, store)

	s := server.NewRedisServer(transport.NewTcpTransport(":9871"))
	if err := s.Serve(commandsServer.Handlers()); err != nil {
		panic(err)
	}
}
