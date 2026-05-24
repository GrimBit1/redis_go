package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Test with official client
	rdb := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		Password:    "", // no password
		DB:          0,  // use default DB
		Protocol:    2,
		DialTimeout: time.Minute,
		ReadTimeout: time.Minute,
	})

	fmt.Println("connected")
	ctx := context.Background()

	err := rdb.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {

		panic(err)
	}

	val, err := rdb.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val) // >>> foo bar
}
