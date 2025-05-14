package cmds

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

// StatusCommand represents the status command
type StatusCommand struct {
	BaseCommand
}

// NewStatusCommand creates a new status command
func NewStatusCommand() *StatusCommand {
	return &StatusCommand{
		BaseCommand: BaseCommand{
			Name:        "status",
			Description: "Check the system status",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	}
}

// Execute handles the command execution
func (c *StatusCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate, logger *slog.Logger) error {
	// TODO: Implement actual status check
	// system.System.ExecStatus()

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "System status: OK",
		},
	})
}
