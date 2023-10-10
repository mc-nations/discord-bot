package handlers

import (
	"fmt"
	"nations/config"
	"strings"

	"nations/discord"
	"nations/redis"

	"github.com/bwmarrin/discordgo"
)

func ListenToServerEvents() {
	redisClient, _ := redis.NewRedisClient()
	mc_server := redisClient.Subscribe("mc_server")
	bot, _ := discord.NewDiscordClient()

	mc_server.RegisterListener("server_start", func(data redis.Json) {
		id := config.GetStr("channel")

		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Title:       "Server Started",
			Description: "The server has just started and you can join now!",
			Color:       0x00ff00, // Green color
		}

		// Send the embed message to the specified channel
		bot.Client.ChannelMessageSendEmbed(id, embed)
	})

	mc_server.RegisterListener("server_lock", func(data redis.Json) {
		id := config.GetStr("channel")

		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Title:       "Server Closed",
			Description: "The server is now closed, it will be open at 'insert time here <@188659395540811776>'",
			Color:       0xff0000, // Red color
		}

		// Send the embed message to the specified channel
		bot.Client.ChannelMessageSendEmbed(id, embed)
	})

	mc_server.RegisterListener("server_unlock", func(data redis.Json) {
		id := config.GetStr("channel")

		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Title:       "Server Opened",
			Description: "The server is now open and you can join now.",
			Color:       0x00ff00, // Green color
		}

		// Send the embed message to the specified channel
		bot.Client.ChannelMessageSendEmbed(id, embed)
	})
}

func ListenToPlayerEvents() {
	redisClient, _ := redis.NewRedisClient()
	mc_server := redisClient.Subscribe("mc_server")
	bot, _ := discord.NewDiscordClient()

	mc_server.RegisterListener("player_join", func(data redis.Json) {
		channel_id := config.GetStr("channel")
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
			Color:       0x00ff00, // Green color
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://crafatar.com/avatars/" + strings.Replace(minecraft_id, "-", "", -1) + ".png?size=128",
				Name:    minecraft_name,
			},
		}

		// Send the embed message to the specified channel
		_, err = bot.Client.ChannelMessageSendEmbed(channel_id, embed)
		if err != nil {
			fmt.Println(err)
		}

	})

	mc_server.RegisterListener("player_quit", func(data redis.Json) {
		channel_id := config.GetStr("channel")
		minecraft_id := data["minecraft_user"].(redis.Json)["id"].(string)
		minecraft_name := data["minecraft_user"].(redis.Json)["name"].(string)
		description := minecraft_name + " left the server!"
		discordUser := data["discord_user"].(redis.Json)
		fmt.Println(discordUser)
		if discordUser != nil && discordUser["id"] != nil {
			member_id := discordUser["id"].(string)
			user, err := bot.Client.User(member_id)
			if err != nil {
				fmt.Println(err)
			}
			description = user.Mention() + " left the server!"
		}

		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Description: description,
			Color:       0xff0000, // Green color
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://crafatar.com/avatars/" + strings.Replace(minecraft_id, "-", "", -1) + ".png?size=128",
				Name:    minecraft_name,
			},
		}

		// Send the embed message to the specified channel
		_, err := bot.Client.ChannelMessageSendEmbed(channel_id, embed)
		if err != nil {
			fmt.Println(err)
		}

	})
}
