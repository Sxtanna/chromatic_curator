package discord

import (
	"bytes"
	"emperror.dev/errors"
	"github.com/Sxtanna/chromatic_curator/internal/app/discord/cmds"
	"github.com/Sxtanna/chromatic_curator/internal/app/discord/data"
	"github.com/Sxtanna/chromatic_curator/internal/common"
	"github.com/Sxtanna/chromatic_curator/internal/system/backend"
	discord "github.com/bwmarrin/discordgo"
	"log/slog"
	"strings"
)

const (
	configurationMissing = errors.Sentinel("discord configuration missing")
)

type BotService struct {
	Bot    *discord.Session
	Config *BotConfiguration

	Logger  *slog.Logger
	Backend backend.Backend

	commands       *cmds.Registry
	registeredCmds map[string][]*discord.ApplicationCommand
}

func (d *BotService) Init(config common.Configuration) error {
	discordConfiguration := common.FindConfiguration[BotConfiguration](config)
	if discordConfiguration == nil {
		return configurationMissing
	}

	d.registeredCmds = make(map[string][]*discord.ApplicationCommand)

	session, err := discord.New("Bot " + discordConfiguration.Token)
	if err != nil {
		return errors.Wrap(err, "failed to create discord session")
	}

	d.Bot = session
	d.Config = discordConfiguration

	// Initialize command registry
	d.commands = cmds.NewRegistry(d.Logger)

	// Register commands
	d.commands.RegisterCommand(cmds.NewRoleCommand(d.Backend, func(id string) bool { return strings.Contains(d.Config.Admins, id) }))
	d.commands.RegisterCommand(cmds.NewColorCommand())

	return nil
}

func (d *BotService) Start() error {
	d.Bot.Identify.Intents = discord.IntentsAll

	d.Bot.AddHandlerOnce(func(s *discord.Session, event *discord.Disconnect) {
		d.Logger.Info("Discord Session has been disconnected!")
	})

	// Add handler for slash commands
	d.Bot.AddHandler(func(s *discord.Session, i *discord.InteractionCreate) {
		if i.Type != discord.InteractionApplicationCommand {
			return
		}

		commandName := i.ApplicationCommandData().Name
		d.Logger.Info("Received slash command",
			slog.String("command", commandName))

		// Get the command from the registry
		cmd, exists := d.commands.GetCommand(commandName)
		if !exists {
			err := s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
				Type: discord.InteractionResponseChannelMessageWithSource,
				Data: &discord.InteractionResponseData{
					Content: "Unknown command: " + commandName,
					Flags:   discord.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				d.Logger.Error("Failed to respond to unknown command",
					slog.String("command", commandName),
					slog.Any("error", err))
			}
			return
		}

		// Execute the command
		err := cmd.Execute(s, i, d.Logger)
		if err != nil {
			d.Logger.Error("Failed to execute command",
				slog.String("command", commandName),
				slog.Any("error", err))
		}
	})

	d.Bot.AddHandler(func(s *discord.Session, i *discord.InteractionCreate) {
		// Only handle button interactions
		if i.Type != discord.InteractionMessageComponent {
			return
		}

		// Get the custom ID from the interaction
		customID := i.MessageComponentData().CustomID

		// Check if this is a share color button
		if !strings.HasPrefix(customID, "share_color:") {
			return
		}

		d.Logger.Info("Received share color button click",
			slog.String("customID", customID))

		// Parse the custom ID to get the image ID
		// Format: share_color:<uuid>
		parts := strings.Split(customID, ":")
		if len(parts) != 2 {
			d.Logger.Error("Invalid custom ID format",
				slog.String("customID", customID))
			return
		}

		imageID := parts[1]

		// Retrieve the image from the cache
		generation := data.FindGeneration(imageID)

		if generation == nil {
			d.Logger.Error("Image not found in cache",
				slog.String("imageID", imageID))
			return
		}

		// Acknowledge the interaction
		err := s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
			Type: discord.InteractionResponseChannelMessageWithSource,
			Data: &discord.InteractionResponseData{
				Embeds: []*discord.MessageEmbed{generation.Embed},
				Files: []*discord.File{
					{
						Name:   "color_preview.png",
						Reader: bytes.NewReader(generation.ImageData),
					},
				},
			},
		})

		if err != nil {
			d.Logger.Error("Failed to respond to button interaction",
				slog.Any("error", err))

			empty := make([]discord.MessageComponent, 0)

			_, err = s.FollowupMessageEdit(i.Interaction, generation.TempMsgID, &discord.WebhookEdit{Components: &empty})

			d.Logger.Error("Failed to remove button interaction",
				slog.Any("error", err))

			return
		}

		err = s.FollowupMessageDelete(i.Interaction, generation.TempMsgID)
		if err != nil {
			d.Logger.Error("Failed to delete temp message",
				slog.Any("error", err))
		}
	})

	if err := d.Bot.Open(); err != nil {
		return errors.Wrap(err, "failed to open bot session")
	}

	d.Logger.Debug("bot session has been opened, registering commands...")

	// Register global commands with Discord
	if err := d.registerCommands(""); err != nil {
		return errors.Wrap(err, "failed to register commands")
	}

	d.Logger.Debug("commands registered, service start complete...")

	return common.ServiceStartedNormallyButDoesNotBlock
}

func (d *BotService) Close(_ error) error {
	d.Logger.Debug("bot close requested, enabling sync events...")
	d.Bot.SyncEvents = true

	// delete global commands
	// d.deleteRegisteredCommands("")

	return d.Bot.Close()
}

func (d *BotService) registerCommands(guildID string) error {
	// Get all commands from the registry
	commands := d.commands.GetApplicationCommands()

	d.Logger.Info("Registering slash commands",
		slog.Int("count", len(commands)),
		slog.String("guild", guildID))

	d.registeredCmds[guildID] = make([]*discord.ApplicationCommand, 0)

	// Register commands with Discord
	for _, cmd := range commands {
		registeredCmd, err := d.Bot.ApplicationCommandCreate(d.Bot.State.User.ID, guildID, cmd)
		if err != nil {
			d.Logger.Error("failed to register slash command",
				slog.String("command", cmd.Name),
				slog.Any("error", err))

			return errors.Wrap(err, "failed to register slash command")
		}

		d.registeredCmds[guildID] = append(d.registeredCmds[guildID], registeredCmd)

		d.Logger.Info("registered slash command",
			slog.String("command", cmd.Name))
	}

	return nil
}

func (d *BotService) deleteRegisteredCommands(guildID string) {
	d.Logger.Info("deleting registered slash commands",
		slog.String("guild", guildID))

	deleteRegisteredCommand := func(guildID string, cmd *discord.ApplicationCommand) {
		err := d.Bot.ApplicationCommandDelete(d.Bot.State.User.ID, guildID, cmd.ID)
		if err != nil {
			d.Logger.Error("failed to delete slash command",
				slog.String("command", cmd.Name),
				slog.Any("error", err))
		} else {
			d.Logger.Info("deleted slash command",
				slog.String("command", cmd.Name))
		}
	}

	// If we have tracked registered commands, use those
	if d.registeredCmds[guildID] != nil && len(d.registeredCmds[guildID]) > 0 {
		registeredCmds := d.registeredCmds[guildID]

		d.Logger.Info("deleting tracked registered commands",
			slog.Int("count", len(registeredCmds)))

		for _, cmd := range registeredCmds {
			deleteRegisteredCommand(guildID, cmd)
		}

		// Clear the registered commands list
		delete(d.registeredCmds, guildID)
		return
	}

	// Fallback: fetch and delete all commands
	commands, err := d.Bot.ApplicationCommands(d.Bot.State.User.ID, guildID)
	if err != nil {
		d.Logger.Error("failed to fetch slash commands",
			slog.Any("error", err))
		return
	}

	d.Logger.Info("deleting fetched slash commands",
		slog.Int("count", len(commands)))

	for _, cmd := range commands {
		deleteRegisteredCommand(guildID, cmd)
	}
}
