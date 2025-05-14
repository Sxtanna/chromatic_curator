package cmds

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

// EchoCommand represents the echo command
type EchoCommand struct {
	BaseCommand
}

// NewEchoCommand creates a new echo command
func NewEchoCommand() *EchoCommand {
	return &EchoCommand{
		BaseCommand: BaseCommand{
			Name:        "echo",
			Description: "Echo back a message",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The message to echo back",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "ephemeral",
					Description: "Whether the response should be ephemeral (only visible to you)",
					Required:    false,
				},
			},
		},
	}
}

// Execute handles the command execution
func (c *EchoCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate, logger *slog.Logger) error {
	options := i.ApplicationCommandData().Options

	// Get the message option
	var message string
	var ephemeral bool

	for _, opt := range options {
		switch opt.Name {
		case "message":
			message = opt.StringValue()
		case "ephemeral":
			ephemeral = opt.BoolValue()
		}
	}

	logger.Info("Executing echo command", map[string]interface{}{
		"message":   message,
		"ephemeral": ephemeral,
	})

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags: func() discordgo.MessageFlags {
				if ephemeral {
					return discordgo.MessageFlagsEphemeral
				} else {
					return 0
				}
			}(),
		},
	})
}
