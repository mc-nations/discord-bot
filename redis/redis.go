package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"nations/utils"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	client   *redis.Client
	channels map[string]RedisChannel
}

var redisClient *RedisClient = nil

func NewRedisClient() (*RedisClient, error) {

	if redisClient != nil {
		return redisClient, nil
	}
	var client *redis.Client = nil
	err := utils.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Username: os.Getenv("REDIS_USER"),
			Password: os.Getenv("REDIS_PASSWORD"),
		})
		return nil
	}, 12, 10*time.Second)
	if err != nil {
		return nil, err
	}
	redisClient = &RedisClient{client, make(map[string]RedisChannel)}
	fmt.Println("Connected to Redis")
	return redisClient, nil
}

func (r RedisClient) GetChannels() map[string]RedisChannel {
	return r.channels
}

func (r RedisClient) Subscribe(channel string) RedisChannel {
	_, channel_exists := r.channels[channel]
	if channel_exists {
		return r.channels[channel]
	}
	pubsub := r.client.Subscribe(ctx, channel)
	fmt.Println("Subscribed to " + channel)
	r.channels[channel] = RedisChannel{pubsub, make(map[string][]func(Json))}
	return r.channels[channel]
}

type RedisChannel struct {
	pubsub          *redis.PubSub
	event_listeners map[string][]func(Json)
}

type Message struct {
	Type    string `json:"type"`
	Ip      string `json:"ip"`
	Content Json   `json:"content"`
}

type Json = map[string]interface{}

func (c RedisChannel) RegisterListener(name string, callback func(Json)) {
	c.event_listeners[name] = append(c.event_listeners[name], callback)
}

func (c RedisChannel) StartListing() {
	for {
		msg, err := c.pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		var message Message
		json.Unmarshal([]byte(msg.Payload), &message)
		for _, callback := range c.event_listeners[message.Content["name"].(string)] {
			go callback(message.Content)
		}
	}
}
