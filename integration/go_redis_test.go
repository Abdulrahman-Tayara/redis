package integration

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"redis/internal/commands"
	"redis/internal/configs"
	"redis/internal/server"
	store2 "redis/internal/store"
	"redis/pkg/ds"
	"redis/pkg/transport"
	"testing"
	"time"
)

func runServer(port string) *server.RedisServer {
	hashTable := ds.NewExpiringHashTable()

	store := store2.NewStore(hashTable)

	cfg := &configs.Configs{
		Version:      "6.0.3",
		ProtoVersion: 3,
		Mode:         "standalone",
		Modules:      []string{},
		Port:         port,
	}

	commandsServer := commands.NewServer(cfg, store)

	s := server.NewRedisServer(transport.NewTcpTransport(cfg.Address()))

	go func() {
		if err := s.Serve(commandsServer.Handlers()); err != nil {
			panic(err)
		}
	}()

	return s
}

func TestGoRedis(t *testing.T) {
	port := "9871"

	s := runServer(port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()
	defer s.Close()

	ctx := context.TODO()

	testSetAndGet(ctx, rdb)

	testSetAndGetExpiredKey(ctx, rdb)
}

func testSetAndGet(ctx context.Context, rdb *redis.Client) {
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	v, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	log.Println("key", v)
}

func testSetAndGetExpiredKey(ctx context.Context, rdb *redis.Client) {
	err := rdb.Set(ctx, "key2", "value", time.Second*2).Err()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 3)

	_, err = rdb.Get(ctx, "key2").Result()
	if errors.Is(err, redis.Nil) {
		return
	}

	if err == nil {
		panic(errors.New("key should be expired"))
	}
	panic(err)
}
