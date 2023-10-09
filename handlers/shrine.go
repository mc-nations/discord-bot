package handlers

import (
	"fmt"
	"nations/discord"
	"nations/redis"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ListenToShrineEvents listens to events related to the shrine in the Minecraft server and sends messages to Discord users accordingly.
// It subscribes to Redis channels and registers listeners for "player_died", "shrine_player_token_picked_up", and "shrine_received_player_token" events.
// When a player dies, it sends a message to the Discord user informing them that another player needs to pick up their totem to revive them.
// When a player picks up another player's token, it sends a message to the token owner informing them who picked up their token and that they will be revived after a certain amount of time.
// When a player brings a token to the shrine, it sends a message to the token owner informing them that their token has been brought to the shrine and they will be able to join the server again after a certain amount of time.

func ListenToShrineEvents() {
	redisClient, _ := redis.NewRedisClient()
	mc_server := redisClient.Subscribe("mc_server")
	bot, _ := discord.NewDiscordClient()
	mc_server.RegisterListener("player_died", func(data redis.Json) {
		userID := data["discord_user"].(redis.Json)["id"].(string)

		channel, err := bot.Client.UserChannelCreate(userID)
		if err != nil {
			fmt.Println("Error creating DM channel:", err)
			return
		}

		embed := &discordgo.MessageEmbed{
			Title:       "You died!",
			Description: "A player has to pick up your totem in order to revive you!",
			Color:       0x00ff00, // Green color
		}

		_, err = bot.Client.ChannelMessageSendEmbed(channel.ID, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	})

	mc_server.RegisterListener("shrine_player_token_picked_up", func(data redis.Json) {
		actionUser := data["action_user"].(redis.Json)
		tokenUser := data["token_user"].(redis.Json)
		actionUserID := actionUser["discord_user"].(redis.Json)["id"].(string)
		tokenUserID := tokenUser["discord_user"].(redis.Json)["id"].(string)

		channel, err := bot.Client.UserChannelCreate(tokenUserID)
		if err != nil {
			fmt.Println("Error creating DM channel:", err)
			return
		}

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://crafatar.com/avatars/" + strings.Replace(actionUser["minecraft_user"].(redis.Json)["id"].(string), "-", "", -1) + ".png?size=128",
				Name:    "Your token has been picked up!",
			},
			Description: fmt.Sprintf("Your token has been picked up by <@%s>. Once they bring it to the shrine you will be revived after xx hours.", actionUserID),
			Color:       0x6c0094, // Purple color
		}

		_, err = bot.Client.ChannelMessageSendEmbed(channel.ID, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	})

	mc_server.RegisterListener("shrine_received_player_token", func(data redis.Json) {
		tokenUserID := data["discord_user"].(redis.Json)["id"].(string)

		channel, err := bot.Client.UserChannelCreate(tokenUserID)
		if err != nil {
			fmt.Println("Error creating DM channel:", err)
			return
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Your token has been brought to the shrine!",
			Description: "You will be able to join the server again after 24 hours.",
			Color:       0x6c0094, // Purple color
		}

		_, err = bot.Client.ChannelMessageSendEmbed(channel.ID, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	})
}
