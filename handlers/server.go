package handlers

import (
	"fmt"
	"nations/config"
	"slices"
	"strings"

	"nations/discord"
	"nations/redis"

	"github.com/bwmarrin/discordgo"
)

func ListenToServerEvents() {
	redisClient, _ := redis.NewRedisClient()
	mc_server := redisClient.Subscribe("mc_server")
	bot, _ := discord.NewDiscordClient()

	/* mc_server.RegisterListener("server_start", func(data redis.Json) {
		id := config.GetStr("channel")

		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Title:       "Server Started",
			Description: "The server has just started and you can join now!",
			Color:       0x00ff00, // Green color
		}

		// Send the embed message to the specified channel
		bot.Client.ChannelMessageSendEmbed(id, embed)
	}) */

	mc_server.RegisterListener("server_lock", func(data redis.Json) {
		id := config.GetStr("channel")

		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Title:       "Server geschlossen",
			Description: "Der Server ist nun geschlossen!",
			Color:       0xff0000, // Red color
		}

		// Send the embed message to the specified channel
		bot.Client.ChannelMessageSendEmbed(id, embed)

		removeOnlineRoles(bot)
	})

	mc_server.RegisterListener("server_unlock", func(data redis.Json) {
		id := config.GetStr("channel")

		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Title:       "Server geöffnet",
			Description: "Der Server ist nun geöffnet, du kannst nun joinen!",
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
			fmt.Println("error getting user (join event)")
			fmt.Println(err)
		}
		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Description: user.Mention() + " hat den Server betreten!",
			Color:       0x00ff00, // Green color
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://crafatar.com/avatars/" + strings.Replace(minecraft_id, "-", "", -1) + ".png?size=128",
				Name:    minecraft_name,
			},
		}

		// Send the embed message to the specified channel
		_, err = bot.Client.ChannelMessageSendEmbed(channel_id, embed)
		if err != nil {
			fmt.Println("error sending join message")
			fmt.Println(err)
		}

		guild_id := config.GetStr("guild")
		role_id := config.GetStr("onlineRole")
		roleErr := bot.Client.GuildMemberRoleAdd(guild_id, member_id, role_id)
		if roleErr != nil {
			fmt.Println("error adding role")
			fmt.Println(roleErr)
		}
	})

	mc_server.RegisterListener("player_quit", func(data redis.Json) {
		channel_id := config.GetStr("channel")
		member_id := data["discord_user"].(redis.Json)["id"].(string)
		minecraft_id := data["minecraft_user"].(redis.Json)["id"].(string)
		minecraft_name := data["minecraft_user"].(redis.Json)["name"].(string)
		description := minecraft_name + " hat den Server verlassen!"
		if data["discord_user"] != nil {
			discordUser := data["discord_user"].(redis.Json)
			if discordUser != nil && discordUser["id"] != nil {
				member_id := discordUser["id"].(string)
				user, err := bot.Client.User(member_id)
				if err != nil {
					fmt.Println("error getting discord user")
					fmt.Println(err)
				}
				description = user.Mention() + " hat den Server verlassen!"
			}
		}

		// Create the embed message
		embed := &discordgo.MessageEmbed{
			Description: description,
			Color:       0xf5d142, // yellow color
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://crafatar.com/avatars/" + strings.Replace(minecraft_id, "-", "", -1) + ".png?size=128",
				Name:    minecraft_name,
			},
		}

		// Send the embed message to the specified channel
		_, err := bot.Client.ChannelMessageSendEmbed(channel_id, embed)
		if err != nil {
			fmt.Println("error sending quit message")
			fmt.Println(err)
		}

		guild_id := config.GetStr("guild")
		role_id := config.GetStr("onlineRole")

		roleErr := bot.Client.GuildMemberRoleRemove(guild_id, member_id, role_id)
		if err != nil {
			fmt.Println("error removing role")
			fmt.Println(roleErr)
		}
	})
}

func removeOnlineRoles(bot *discord.DiscordBot) {
	guild_id := config.GetStr("guild")
	role_id := config.GetStr("onlineRole")

	members, err := bot.Client.GuildMembers(guild_id, "", 1000)
	if err != nil {
		fmt.Print("error while getting guild members")
		fmt.Print(err)
	}
	for _, member := range members {
		if slices.Contains(member.Roles, role_id) {
			roleErr := bot.Client.GuildMemberRoleRemove(guild_id, member.User.ID, role_id)
			if err != nil {
				fmt.Println("error removing role")
				fmt.Println(roleErr)
			}
		}
	}
}
