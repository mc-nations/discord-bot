package handlers

import (
	"fmt"
	"nations/config"
	"nations/discord"
	"nations/redis"
)

// ListenToAccountLinkEvents listens to Redis events for account linking and updates Discord roles accordingly.
func ListenToAccountLinkEvents() {
	redisClient, _ := redis.NewRedisClient()
	mc_server := redisClient.Subscribe("mc_server")
	bot, _ := discord.NewDiscordClient()

	// Register listener for when a Minecraft account is linked to a Discord account
	mc_server.RegisterListener("discord_account_linked", func(data redis.Json) {
		guild_id := config.GetStr("guild")
		member_id := data["discord_user"].(redis.Json)["id"].(string)
		role_id := config.GetStr("linkedRole")
		err := bot.Client.GuildMemberRoleAdd(guild_id, member_id, role_id)
		if err != nil {
			fmt.Println("error adding role")
			fmt.Println(err)
		}
	})

	// Register listener for when a Minecraft account is unlinked from a Discord account
	mc_server.RegisterListener("discord_account_unlinked", func(data redis.Json) {
		guild_id := config.GetStr("guild")
		member_id := data["discord_user"].(redis.Json)["id"].(string)
		role_id := config.GetStr("linkedRole")
		err := bot.Client.GuildMemberRoleRemove(guild_id, member_id, role_id)
		if err != nil {
			fmt.Println("error removing role")
			fmt.Println(err)
		}
	})
}
