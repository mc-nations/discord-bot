package discord

import (
	"fmt"
	"log"
	"nations/config"
	"nations/redis"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	adminPermission int64 = discordgo.PermissionManageServer
	commands              = []*discordgo.ApplicationCommand{
		{
			Name:                     "here",
			Description:              "Set the channel to send notifications to",
			DefaultMemberPermissions: &adminPermission,
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"here": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			id := i.Interaction.ChannelID
			config.Save("channel", id)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Now sending all messages to this channel",
				},
			})
		},
	}
)

type DiscordBot struct {
	Client *discordgo.Session
}

func NewDiscordClient(client redis.RedisClient) DiscordBot {
	s, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic("error login to discord")
	}
	err = s.Open()
	if err != nil {
		panic("error opening connection")
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		fmt.Println(i.ApplicationCommandData().Name)
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	//dafuq
	// invalid memory address or nil pointer derefrence
	s.Identify.Intents = discordgo.IntentGuildMessages
	s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", []*discordgo.ApplicationCommand{})
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {

		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {

			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	return DiscordBot{Client: s}
}
