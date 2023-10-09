package main

import (
	"fmt"
	"nations/handlers"
	"nations/redis"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	client, err := redis.NewRedisClient()
	if err != nil {
		panic("failed to login to redis")
	}

	handlers.ListenToPlayerEvents()
	handlers.ListenToAccountLinkEvents()
	handlers.ListenToShrineEvents()

	mc_server := client.Subscribe("mc_server")

	mc_server.RegisterListener("server_remaining_uptime", func(data redis.Json) {
		fmt.Println(data["remaining_time"])
	})

	mc_server.StartListing()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
