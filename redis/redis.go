package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient(config redis.Options) RedisClient {
	return RedisClient{redis.NewClient(&config)}
}

type Subscriber interface {
	Subscribe(channel string) RedisChannel
}
type RedisClient struct {
	client *redis.Client
}

func (r RedisClient) Subscribe(channel string) RedisChannel {
	fmt.Println("Subscribing to channel", channel)
	pubsub := r.client.Subscribe(ctx, channel)
	return RedisChannel{pubsub}
}

type Listener interface {
	On(name string, callback func(data interface{}))
}
type RedisChannel struct {
	pubsub *redis.PubSub
}

type Content struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Message struct {
	Type    string  `json:"type"`
	Ip      string  `json:"ip"`
	Content Content `json:"content"`
}

func (c RedisChannel) On(name string, callback func(data interface{})) {
	go func() {
		for {
			raw, err := c.pubsub.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}
			var msg Message
			err = json.Unmarshal([]byte(raw.Payload), &msg)
			if err != nil {
				panic(err)
			}
			if msg.Content.Name == name {
				callback(msg.Content.Data)
			}
		}
	}()
}
