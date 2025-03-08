package main

import (
	"context"
	"flag"
	"os/signal"
	"redis/internal/commands"
	"redis/internal/configs"
	"redis/internal/server"
	store2 "redis/internal/store"
	"redis/logger"
	"redis/pkg/active_expiration"
	"redis/pkg/ds"
	"redis/pkg/transport"
	"syscall"
)

var (
	configFilePath = ""
)

func init() {
	flag.StringVar(&configFilePath, "config", "config.json", "config file path")
	flag.Parse()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	cfg, err := configs.LoadConfigsOrDefaults(configFilePath)
	if err != nil {
		panic(err)
	}

	hashTable := ds.NewExpiringHashTable()

	store := store2.NewStore(hashTable)

	runActiveExpirationLoop(cfg, hashTable)

	commandsServer := commands.NewServer(cfg, store)

	s := server.NewRedisServer(transport.NewTcpTransport(cfg.Address()))

	logger.Infof("Starting server on %s...", cfg.Port)

	go func() {
		if err := s.Serve(commandsServer.Handlers()); err != nil {
			panic(err)
		}
	}()

	// listen for the interrupt signal
	<-ctx.Done()

	logger.Infof("Shutting down server on %s...", cfg.Port)

	// stop the serve
	if err = s.Close(); err != nil {
		panic(err)
	}
}

func runActiveExpirationLoop(_ *configs.Configs, hashTable *ds.ExpiringHashTable) {
	ae := active_expiration.NewRandomActiveExpiration(hashTable)

	go ae.Start()
}
