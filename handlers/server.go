package handlers

import (
	"fmt"
	"strings"

	"nations/redis"
	"nations/discord"
	"github.com/bwmarrin/discordgo"
)

func ListenToPlayerEvents() {
	redisClient, _ := redis.NewRedisClient()
	mc_server := redisClient.Subscribe("mc_server")
	bot, _ := discord.NewDiscordClient()
	mc_server.RegisterListener("player_join", func(data redis.Json) {
		channel_id := "your_channel_id" // replace with your channel ID
		member_id := data["discord_user"].(redis.Json)["id"].(string)
		minecraft_id := data["minecraft_user"].(redis.Json)["id"].(string)
		minecraft_name := data["minecraft_user"].(redis.Json)["name"].(string)
		user, err := bot.Client.User(member_id)
		if err != nil {
			fmt.Println(err)
		}
		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Description: user.Mention() + " joined the server!",
			Color: 0x00ff00, // Green color
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://crafatar.com/avatars/" + strings.Replace(minecraft_id, "-", "", -1) + ".png?size=128",
				Name: minecraft_name,
			},
		}

		// Send the embed message to the specified channel
		_, err = bot.Client.ChannelMessageSendEmbed(channel_id, embed)
		if err != nil {
			fmt.Println(err)
		}

	
	})

	mc_server.RegisterListener("player_leave", func(data redis.Json) {
		channel_id := "your_channel_id" // replace with your channel ID
		member_id := data["discord_user"].(redis.Json)["id"].(string)
		minecraft_id := data["minecraft_user"].(redis.Json)["id"].(string)
		minecraft_name := data["minecraft_user"].(redis.Json)["name"].(string)
		user, err := bot.Client.User(member_id)
		if err != nil {
			fmt.Println(err)
		}
		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Description: user.Mention() + " left the server!",
			Color: 0xff0000, // Green color
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://crafatar.com/avatars/" + strings.Replace(minecraft_id, "-", "", -1) + ".png?size=128",
				Name: minecraft_name,
			},
		}

		// Send the embed message to the specified channel
		_, err = bot.Client.ChannelMessageSendEmbed(channel_id, embed)
		if err != nil {
			fmt.Println(err)
		}

	
	})
}