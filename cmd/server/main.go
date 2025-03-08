package main

import (
	"redis/internal/commands"
	"redis/internal/configs"
	"redis/internal/server"
	store2 "redis/internal/store"
	"redis/logger"
	"redis/pkg/active_expiration"
	"redis/pkg/ds"
	"redis/pkg/transport"
)

func main() {

	hashTable := ds.NewExpiringHashTable()

	store := store2.NewStore(hashTable)

	cfg := &configs.Configs{
		Version:      "6.0.3",
		ProtoVersion: 3,
		Mode:         "standalone",
		Modules:      []string{},
		Port:         "9871",
	}

	runActiveExpirationLoop(cfg, hashTable)

	commandsServer := commands.NewServer(cfg, store)

	s := server.NewRedisServer(transport.NewTcpTransport(cfg.Address()))

	logger.Infof("Starting server on %s...", cfg.Port)

	if err := s.Serve(commandsServer.Handlers()); err != nil {
		panic(err)
	}
}

func runActiveExpirationLoop(_ *configs.Configs, hashTable *ds.ExpiringHashTable) {
	ae := active_expiration.NewRandomActiveExpiration(hashTable)

	go ae.Start()
}
