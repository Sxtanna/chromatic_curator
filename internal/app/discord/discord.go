package discord

import (
	"emperror.dev/errors"
	"github.com/Sxtanna/chromatic_curator/internal/common"
	discord "github.com/bwmarrin/discordgo"
	"log/slog"
)

const (
	configurationMissing = errors.Sentinel("discord configuration missing")
)

type BotService struct {
	Bot    *discord.Session
	Config *BotConfiguration

	Logger *slog.Logger
}

func (d *BotService) Init(config common.Configuration) error {
	discordConfiguration := common.FindConfiguration[BotConfiguration](config)
	if discordConfiguration == nil {
		return configurationMissing
	}

	session, err := discord.New("Bot " + discordConfiguration.Token)
	if err != nil {
		return errors.Wrap(err, "failed to create discord session")
	}

	d.Bot = session
	d.Config = discordConfiguration

	return nil
}

func (d *BotService) Start() error {
	d.Bot.Identify.Intents = discord.IntentsAll

	// Define slash commands
	commands := []*discord.ApplicationCommand{
		{
			Name:        "status",
			Description: "Check the system status",
		},
		{
			Name:        "users",
			Description: "Get the list of users",
		},
	}

	d.Bot.AddHandlerOnce(func(s *discord.Session, event *discord.Connect) {
		d.Logger.Info("Discord Session is now ready!")

		// Register slash commands
		guildID := ""
		if d.Config != nil && d.Config.GuildID != "" {
			guildID = d.Config.GuildID
		}

		for _, cmd := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
			if err != nil {
				d.Logger.Error("Failed to create slash command", map[string]interface{}{
					"command": cmd.Name,
					"error":   err.Error(),
				})
			} else {
				d.Logger.Info("Registered slash command", map[string]interface{}{
					"command": cmd.Name,
					"guild":   guildID,
				})
			}
		}
	})

	d.Bot.AddHandlerOnce(func(_ *discord.Session, event *discord.Disconnect) {
		d.Logger.Info("Discord Session has been disconnected!")
	})

	// Add handler for slash commands
	d.Bot.AddHandler(func(s *discord.Session, i *discord.InteractionCreate) {
		if i.Type != discord.InteractionApplicationCommand {
			return
		}

		d.Logger.Info("Received slash command", map[string]interface{}{
			"command": i.ApplicationCommandData().Name,
		})

		switch i.ApplicationCommandData().Name {
		case "status":
			// TODO: Implement actual status check
			// system.System.ExecStatus()
			err := s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
				Type: discord.InteractionResponseChannelMessageWithSource,
				Data: &discord.InteractionResponseData{
					Content: "System status: OK",
				},
			})
			if err != nil {
				d.Logger.Error("Failed to respond to interaction", map[string]interface{}{
					"command": i.ApplicationCommandData().Name,
					"error":   err.Error(),
				})
			}
		case "users":
			// TODO: Implement actual user list
			// system.System.ExecGetUsers()
			err := s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
				Type: discord.InteractionResponseChannelMessageWithSource,
				Data: &discord.InteractionResponseData{
					Content: "User list would be displayed here",
				},
			})
			if err != nil {
				d.Logger.Error("Failed to respond to interaction", map[string]interface{}{
					"command": i.ApplicationCommandData().Name,
					"error":   err.Error(),
				})
			}
		default:
			err := s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
				Type: discord.InteractionResponseChannelMessageWithSource,
				Data: &discord.InteractionResponseData{
					Content: "Unknown command",
				},
			})
			if err != nil {
				d.Logger.Error("Failed to respond to interaction", map[string]interface{}{
					"command": i.ApplicationCommandData().Name,
					"error":   err.Error(),
				})
			}
		}
	})

	// Keep the message command handler for backward compatibility
	d.Bot.AddHandlerOnce(func(_ *discord.Session, event *discord.MessageCreate) {
		command := strings.TrimPrefix(event.Message.ContentWithMentionsReplaced(), "d!")

		switch command {
		case "status":
			// system.System.ExecStatus()
		case "users":
			// system.System.ExecGetUsers()
		}
	})

	if err := d.Bot.Open(); err != nil {
		return errors.Wrap(err, "failed to open bot session")
	}

	d.Logger.Debug("bot session has been opened, service start complete...")

	return common.ServiceStartedNormallyButDoesNotBlock
}

func (d *BotService) Close(_ error) error {
	d.Logger.Debug("bot close requested, enabling sync events...")
	d.Bot.SyncEvents = true

	return d.Bot.Close()
}
