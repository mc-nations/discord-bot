package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"nations/config"
	"nations/discord"
	"nations/redis"
	"os"
	"strconv"
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
			Title:       "Du bist gestorben!",
			Description: "Ein anderer Spieler muss nun dein Token aufheben, damit du schneller und ohne Strafe wiederbelebt wirst.",
			Color:       0x6c0094, // Purple color
		}

		_, err = bot.Client.ChannelMessageSendEmbed(channel.ID, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		minecraft_id := data["minecraft_user"].(redis.Json)["id"].(string)
		minecraft_name := data["minecraft_user"].(redis.Json)["name"].(string)
		description := minecraft_name + " ist gestorben!"
		if data["discord_user"] != nil {
			discordUser := data["discord_user"].(redis.Json)
			if discordUser != nil && discordUser["id"] != nil {
				member_id := discordUser["id"].(string)
				user, err := bot.Client.User(member_id)
				if err != nil {
					fmt.Println("error getting discord user")
					fmt.Println(err)
				}
				description = user.Mention() + " ist gestorben!"
			}
		}

		global_embed := &discordgo.MessageEmbed{
			Description: description,
			Color:       0xff0000, // red color
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://crafatar.com/avatars/" + strings.Replace(minecraft_id, "-", "", -1) + ".png?size=128",
				Name:    minecraft_name,
			},
		}

		global_channel_id := config.GetStr("channel")
		_, err = bot.Client.ChannelMessageSendEmbed(global_channel_id, global_embed)
		if err != nil {
			fmt.Println("error sending quit message")
			fmt.Println(err)
		}

	})

	mc_server.RegisterListener("shrine_player_token_picked_up", func(data redis.Json) {
		actionUser := data["action_user"].(redis.Json)
		tokenUser := data["token_user"].(redis.Json)
		actionUserID := actionUser["discord_user"].(redis.Json)["id"].(string)
		tokenUserID := tokenUser["discord_user"].(redis.Json)["id"].(string)
		if checkPlayerPickUpCombi(actionUserID, tokenUserID) {
			return
		} else {
			setTokenHolder(actionUserID, tokenUserID)
		}

		channel, err := bot.Client.UserChannelCreate(tokenUserID)
		if err != nil {
			fmt.Println("Error creating DM channel:", err)
			return
		}
		reviveTime := data["revive_time"].(float64)

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://crafatar.com/avatars/" + strings.Replace(actionUser["minecraft_user"].(redis.Json)["id"].(string), "-", "", -1) + ".png?size=128",
				Name:    "Dein Token wurde aufgehoben!",
			},
			Description: fmt.Sprintf("Dein Token wurde von <@%s> aufgehoben, sobald es zum Shrine gebracht wurde wirst du nach %s wiederbelebt.", actionUserID, getReviveTimeString(reviveTime)),
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
		reviveTime := data["revive_time"].(float64)

		embed := &discordgo.MessageEmbed{
			Title:       "Dein Token wurde zum Shrine gebracht!",
			Description: fmt.Sprintf("Du kannst nach %s wieder joinen.", getReviveTimeString(reviveTime)),
			Color:       0x6c0094, // Purple color
		}
		_, err = bot.Client.ChannelMessageSendEmbed(channel.ID, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	})

	mc_server.RegisterListener("shrine_revived_player", func(data redis.Json) {
		userID := data["discord_user"].(redis.Json)["id"].(string)
		reviveType := data["revive_type"].(string)
		channel, err := bot.Client.UserChannelCreate(userID)
		removeTokenHolder(userID)
		if err != nil {
			fmt.Println("Error creating DM channel:", err)
			return
		}

		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Du wurdest%swiederbelebt!", getRevivedByString(reviveType)),
			Description: "Du kannst nun wieder joinen.",
			Color:       0x6c0094, // purple color
		}

		_, err = bot.Client.ChannelMessageSendEmbed(channel.ID, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		minecraft_id := data["minecraft_user"].(redis.Json)["id"].(string)
		minecraft_name := data["minecraft_user"].(redis.Json)["name"].(string)
		description := minecraft_name + " wurde wiederbelebt!"
		if data["discord_user"] != nil {
			discordUser := data["discord_user"].(redis.Json)
			if discordUser != nil && discordUser["id"] != nil {
				member_id := discordUser["id"].(string)
				user, err := bot.Client.User(member_id)
				if err != nil {
					fmt.Println("error getting discord user")
					fmt.Println(err)
				}
				description = user.Mention() + " wurde wiederbelebt!"
			}
		}

		// Create the embed message
		global_embed := &discordgo.MessageEmbed{
			Description: description,
			Color:       0x6c0094, // purple color
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://crafatar.com/avatars/" + strings.Replace(minecraft_id, "-", "", -1) + ".png?size=128",
				Name:    minecraft_name,
			},
		}

		global_channel_id := config.GetStr("channel")
		_, err = bot.Client.ChannelMessageSendEmbed(global_channel_id, global_embed)
		if err != nil {
			fmt.Println("error sending quit message")
			fmt.Println(err)
		}

	})

}
func getRevivedByString(reviveType string) string {
	reviveType = strings.ToLower(reviveType)
	if reviveType == "shrine" {
		return " vom Shrine "
	} else if reviveType == "timer" {
		return " nach Ablauf der Strafzeit "
	} else if reviveType == "command" {
		return " von einem Administrator "
	} else {
		return " "
	}
}

func getReviveTimeString(millis float64) string {
	hours := millis / 1000 / 60 / 60
	rest_minutes := hours*60 - math.Floor(hours)*60
	if hours >= 1 {
		return strconv.FormatFloat(hours, 'f', 0, 64) + " Stunden und " + strconv.FormatFloat(rest_minutes, 'f', 0, 64) + " Minuten"
	} else {
		return strconv.FormatFloat(rest_minutes, 'f', 0, 64) + " Minuten"
	}
}

func readPlayerCombinations() redis.Json {
	if _, err := os.Stat("playerTokenCombinations.json"); err != nil {
		os.WriteFile("playerTokenCombinations.json", []byte("{}"), 0644)
	}
	dat, err := os.ReadFile("playerTokenCombinations.json")
	if err != nil {
		panic("error reading playerTokenCombinations file")
	}
	var combinations redis.Json
	err = json.Unmarshal(dat, &combinations)
	if err != nil {
		panic("error unmarshaling")
	}

	return combinations
}

func checkPlayerPickUpCombi(actionUserID string, tokenUserID string) bool {
	var combinations = readPlayerCombinations()
	holder := combinations[tokenUserID]
	if holder == nil {
		return false
	}
	return holder == actionUserID
}

// addPlayerPickUpCombi adds a combination of a player picking up another player's token to the json file of combinations
// prevent that player receives multiple messages when multiple players pick up their token
func setTokenHolder(actionUserID string, tokenUserID string) {
	var combinations redis.Json
	combinations = readPlayerCombinations()
	combinations[tokenUserID] = actionUserID
	saveFile(combinations)
}

func removeTokenHolder(tokenUserId string) {
	combinations := readPlayerCombinations()
	delete(combinations, tokenUserId)
	saveFile(combinations)
}

// saveFile saves the json file of playerTokenCombinations
func saveFile(combinations redis.Json) {
	dat, err := json.Marshal(combinations)
	if err != nil {
		panic("error marshaling token")
	}
	err = os.WriteFile("playerTokenCombinations.json", dat, 0644)
	if err != nil {
		panic("error writing token")
	}

}
