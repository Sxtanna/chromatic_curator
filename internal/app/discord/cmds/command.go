package cmds

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

// Command represents a Discord slash command
type Command interface {
	// GetName returns the name of the command
	GetName() string

	// GetDescription returns the description of the command
	GetDescription() string

	// GetOptions returns the options/arguments for the command
	GetOptions() []*discordgo.ApplicationCommandOption

	// Execute handles the command execution
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate, logger *slog.Logger) error
}

// BaseCommand provides a basic implementation of the Command interface
type BaseCommand struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
}

// GetName returns the name of the command
func (c *BaseCommand) GetName() string {
	return c.Name
}

// GetDescription returns the description of the command
func (c *BaseCommand) GetDescription() string {
	return c.Description
}

// GetOptions returns the options/arguments for the command
func (c *BaseCommand) GetOptions() []*discordgo.ApplicationCommandOption {
	return c.Options
}

// Registry manages all available commands
type Registry struct {
	commands map[string]Command
	logger   *slog.Logger
}

// NewRegistry creates a new command registry
func NewRegistry(logger *slog.Logger) *Registry {
	return &Registry{
		commands: make(map[string]Command),
		logger:   logger,
	}
}

// RegisterCommand adds a command to the registry
func (r *Registry) RegisterCommand(cmd Command) {
	r.commands[cmd.GetName()] = cmd
	r.logger.Info("Registered command",
		slog.String("command", cmd.GetName()))
}

// GetCommand returns a command by name
func (r *Registry) GetCommand(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

// GetAllCommands returns all registered commands
func (r *Registry) GetAllCommands() []Command {
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// GetApplicationCommands returns all commands as ApplicationCommand objects
func (r *Registry) GetApplicationCommands() []*discordgo.ApplicationCommand {
	cmds := make([]*discordgo.ApplicationCommand, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, &discordgo.ApplicationCommand{
			Name:        cmd.GetName(),
			Description: cmd.GetDescription(),
			Options:     cmd.GetOptions(),
		})
	}
	return cmds
}
