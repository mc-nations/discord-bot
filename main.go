package main

import (
	"nations/config"
	"nations/redis"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"nations/handlers"
	

	_ "github.com/joho/godotenv/autoload"

	"nations/discord"

	"github.com/bwmarrin/discordgo"
)
type DiscordUser struct {
	Id string `json:"id"`
}
type MinecraftUser struct {
	Name string `json:"name"`
	Id string `json:"id"`
}

type AccountLinkContent struct {
	DiscordUser DiscordUser `json:"discord_user"`
	MinecraftUser string `json:"minecraft_user"`
}

type UptimeContent struct {
	RemainingTime interface{}
}



func main() {
	client,  err := redis.NewRedisClient()
	if(err != nil){
		panic("failed to login to redis")
	}

	bot, err := discord.NewDiscordClient()
	if(err != nil){
		panic("failed to login to discord")
	}

	
	handlers.ListenToPlayerEvents()
	mc_server := client.Subscribe("mc_server")

	mc_server.RegisterListener("server_start", func(data redis.Json) {
		id := config.GetStr("channel")
		
		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Title: "Server Started",
			Description: "The server has just started and you can join now!",
			Color: 0x00ff00, // Green color
		}
		
		// Send the embed message to the specified channel
		bot.Client.ChannelMessageSendEmbed(id, embed)
	})

	
	mc_server.RegisterListener("player_died", func(data redis.Json) {
		userID := data["discord_user"].(redis.Json)["id"].(string)

		channel, err := bot.Client.UserChannelCreate(userID)
		if err != nil {
			fmt.Println("Error creating DM channel:", err)
			return
		}

		embed := &discordgo.MessageEmbed{
			Title: "You died!",
			Description: "A player has to pick up your totem in order to revive you!",
			Color: 0x00ff00, // Green color
		}

		_, err = bot.Client.ChannelMessageSendEmbed(channel.ID, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	})
		

	mc_server.RegisterListener("server_remaining_uptime", func(data redis.Json) {
		fmt.Println(data["remaining_time"])
	})

	mc_server.RegisterListener("discord_account_linked", func(data redis.Json) {
		guild_id := config.GetStr("guild")
		member_id := data["discord_user"].(redis.Json)["id"].(string)
		role_id := config.GetStr("linkedRole")
		err := bot.Client.GuildMemberRoleAdd(guild_id, member_id,role_id) 
		if err != nil {
			fmt.Println(err)
		}
	})

	mc_server.RegisterListener("discord_account_unlinked", func(data redis.Json) {
		guild_id := config.GetStr("guild")
		member_id := data["discord_user"].(redis.Json)["id"].(string)
		role_id := config.GetStr("linkedRole")
		err := bot.Client.GuildMemberRoleRemove(guild_id, member_id,role_id) 
		if err != nil {
			fmt.Println(err)
		}
	})

	

	mc_server.StartListing()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}



