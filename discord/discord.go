package discord

import (
	"fmt"
	"nations/config"
	"nations/utils"
	"os"
	"time"

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
			var channel, err = s.Channel(id)
			fmt.Println(channel.Type)
			fmt.Println(discordgo.ChannelTypeGuildText)
			if err != nil || channel.Type != discordgo.ChannelTypeGuildText {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Du kannst diesen Command nur auf einem Server ausführen!",
					},
				})
				return

			}
			config.Save("channel", id)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Alle nachrichten werden nun in diesem Channel gesendet!",
				},
			})
		},
	}
)

type DiscordBot struct {
	Client *discordgo.Session
}

var discordClient *DiscordBot = nil

func NewDiscordClient() (*DiscordBot, error) {

	if discordClient != nil {
		return discordClient, nil
	}
	var s *discordgo.Session = nil
	err := utils.Retry(func() error {
		var err error
		s, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
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
		s.Identify.Intents = discordgo.IntentGuildMessages
		registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
		for i, v := range commands {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
			if err != nil {
				panic("Cannot create comamnd")
			}
			registeredCommands[i] = cmd
		}
		return nil
	}, 6, 20*time.Second)

	if err != nil {
		return nil, err
	}
	fmt.Println("Discord bot is now running.")
	discordClient = &DiscordBot{Client: s}
	return discordClient, nil

}
