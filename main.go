package main

import (
	"nations/config"
	"nations/redis"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"

	"nations/discord"

	raw_red "github.com/redis/go-redis/v9"
)

func main() {

	client := redis.NewRedisClient(raw_red.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	bot := discord.NewDiscordClient(client)
	mc_server := client.Subscribe("mc_server")
	mc_server.On("server_start", func(data interface{}) {
		id := config.Get("channel")
		bot.Client.ChannelMessageSend(id, "Server started")

	})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
