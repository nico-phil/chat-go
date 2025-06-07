package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redisclient struct {
	C *redis.Client
}

func NewRedisClient() Redisclient {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	return Redisclient{C: client}
}

func (rc *Redisclient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := rc.C.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc *Redisclient) Sub(ctx context.Context, userID string) error {
	pubsub := rc.C.Subscribe(ctx, userID)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return err
	}

	go rc.Receive(pubsub)

	return nil
}

func (rc *Redisclient) Pub(ctx context.Context, channel string, message string) error {
	err := rc.C.Publish(ctx, channel, message).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *Redisclient) Receive(pubsub *redis.PubSub) {
	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
		}

		fmt.Println(msg.Channel, msg.Payload)
	}

}
