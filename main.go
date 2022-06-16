package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

func main() {
	done := make(chan struct{})
	redisCtx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go Subscribe(done, redisCtx, rdb, "channel-1")

	exit := time.NewTicker(1 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)
Outer:
	for {
		select {
		case <-ticker.C:
			{
				err := rdb.Publish(redisCtx, "channel-1", "ping").Err()
				if err != nil {
					panic(err)
				}
			}
		case <-exit.C:
			{
				close(done)
				break Outer
			}
		}
	}
}

func Subscribe(done <-chan struct{}, ctx context.Context, rdb *redis.Client, channel string) {
	pubSub := rdb.Subscribe(ctx, channel)
	defer pubSub.Close()
	ch := pubSub.Channel()
	for {
		select {
		case payload := <-ch:
			{
				log.Println(payload)
			}
		case <-done:
			{
				break
			}
		}
	}
}
