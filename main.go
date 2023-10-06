package main

import (
	"nations/config"
	"nations/redis"
	"os"
	"os/signal"
	"syscall"
	"fmt"

	_ "github.com/joho/godotenv/autoload"

	"nations/discord"

	raw_red "github.com/redis/go-redis/v9"
	"github.com/bwmarrin/discordgo"
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
		
		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Title: "Server Started",
			Description: "The server has just started and you can join now!",
			Color: 0x00ff00, // Green color
		}
		
		// Send the embed message to the specified channel
		bot.Client.ChannelMessageSendEmbed(id, embed)
	})
	
	mc_server.On("player_join", func(data interface{}) {

	fmt.Println("hello")

		id := config.Get("channel")
		
		// Get the player's name from the data
		playerName := data.(string)
		
		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Title: "Player Joined",
			Description: fmt.Sprintf("%s has joined the server!", playerName),
			Color: 0x00ff00, // Green color
		}
		
		// Send the embed message to the specified channel
		bot.Client.ChannelMessageSendEmbed(id, embed)
	})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
