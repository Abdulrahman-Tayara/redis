package main

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:9871",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	ctx := context.TODO()

	err := rdb.Set(ctx, "key", "value", time.Second).Err()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 2)

	v, err := rdb.Get(ctx, "key").Result()
	if !errors.Is(err, redis.Nil) {
		panic(err)
	} else if errors.Is(err, redis.Nil) {
		log.Println("not found")
		return
	}

	log.Println("key", v)
}
