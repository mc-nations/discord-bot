package main

import (
	"nations/discord"
	"nations/handlers"
	"nations/redis"
	"os"
	"os/signal"
	"syscall"

	"fmt"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	if len(os.Args) > 1 {
		if os.Args[1] == "distribute-roles" {
			fmt.Println(len(os.Args))
			if len(os.Args) <= 2 {
				fmt.Println("Please provide a csv file with format: userId,biome")
				return
			}
			discord.DistributeRoles(os.Args[2])
			return
		}
	}

	client, err := redis.NewRedisClient()
	if err != nil {
		panic("failed to login to redis")
	}

	handlers.ListenToServerEvents()
	handlers.ListenToPlayerEvents()
	handlers.ListenToAccountLinkEvents()
	handlers.ListenToShrineEvents()

	mc_server := client.Subscribe("mc_server")

	mc_server.StartListing()
	//discord.DistributeRoles()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
