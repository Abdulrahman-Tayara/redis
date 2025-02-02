package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:        "localhost:9871",
		Password:    "", // no password set
		DB:          0,  // use default DB
		ReadTimeout: time.Hour,
		DialTimeout: time.Hour,
	})

	defer rdb.Close()

	ctx := context.TODO()

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
