package cmds

import (
	"context"
	"github.com/Sxtanna/chromatic_curator/internal/common"
	"github.com/Sxtanna/chromatic_curator/internal/system/backend"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"strings"
)

type RoleCommand struct {
	BaseCommand

	backend         backend.Backend
	isAdminFunction func(id string) bool
}

func NewRoleCommand(backend backend.Backend, isAdminFunction func(id string) bool) *RoleCommand {
	return &RoleCommand{
		backend:         backend,
		isAdminFunction: isAdminFunction,
		BaseCommand: BaseCommand{
			Name:        "role",
			Description: "manage your personal role",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "The new name of your role",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "color",
					Description: "The new color of your role",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to modify (yourself by default)",
					Required:    false,
				},
			},
		},
	}
}

func (r *RoleCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate, logger *slog.Logger) error {
	if i.GuildID == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command can only be used in a guild",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	ctx := &RoleUpdateContext{
		ctx:  context.Background(),
		log:  logger,
		bot:  s,
		data: i,
	}

	if g, err := s.Guild(i.GuildID); err == nil {
		ctx.guild = g
	} else {
		logger.Error("failed to get guild",
			slog.Any("error", err),
			slog.String("guild", i.GuildID))

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not get current guild",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	var (
		caller *discordgo.User
		target *discordgo.User
	)

	caller = i.User
	if caller == nil {
		caller = i.Member.User
	}

	target = caller

	if userOption := GetOptionByName(i.Interaction, "user"); userOption != nil {
		specifiedUser := userOption.UserValue(s)

		if specifiedUser == nil {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Could not find user with ID: " + userOption.StringValue() + "",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}

		target = specifiedUser

		if target.ID != caller.ID {
			if !r.isAdminFunction(caller.ID) {
				return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You do not have permission to modify another user's role",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}

			if r.isAdminFunction(target.ID) { // maybe don't keep this
				return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You cannot modify another admin's role",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}
		}
	}

	if target == nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not resolve target user",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	var role *discordgo.Role

	if resolved, err := ctx.resolvePersonalRoleForTarget(caller, target, r); err != nil {
		return err
	} else {
		role = resolved
	}

	responses := make([]string, 0)

	if nameOption := GetOptionByName(i.Interaction, "name"); nameOption != nil {
		err := ctx.updatePersonalRoleName(role, nameOption.StringValue())
		if err == nil {
			responses = append(responses, "Role name updated to \""+nameOption.StringValue()+"\"")
		} else {
			logger.Error("failed to update role name",
				slog.Any("error", err),
				slog.String("role", role.ID),
				slog.String("name", nameOption.StringValue()))

			responses = append(responses, "Could not update role name")
		}
	}

	if colorOption := GetOptionByName(i.Interaction, "color"); colorOption != nil {
		err := ctx.updatePersonalRoleColor(role, colorOption.StringValue())
		if err == nil {
			responses = append(responses, "Role color updated to \""+colorOption.StringValue()+"\"")
		} else {
			logger.Error("failed to update role color",
				slog.Any("error", err),
				slog.String("role", role.ID),
				slog.String("color", colorOption.StringValue()))

			responses = append(responses, "Could not update role color")
		}
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: strings.Join(responses, "\n"),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

type RoleUpdateContext struct {
	ctx context.Context
	log *slog.Logger

	bot   *discordgo.Session
	data  *discordgo.InteractionCreate
	guild *discordgo.Guild
}

func (c *RoleUpdateContext) resolvePersonalRoleForTarget(caller, target *discordgo.User, r *RoleCommand) (*discordgo.Role, error) {
	var role *discordgo.Role

	if existingPersonalRoleID, err := r.backend.GetRole(c.ctx, c.guild.ID, target.ID); err != nil {
		c.log.Error("failed to get role for user",
			slog.Any("error", err),
			slog.String("target", target.ID),
			slog.String("caller", caller.ID))

		return nil, c.bot.InteractionRespond(c.data.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not resolve role for user: " + target.ID,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	} else if existingPersonalRoleID != "" {

		for _, guildRole := range c.guild.Roles {
			if guildRole.ID == existingPersonalRoleID {
				role = guildRole
				break
			}
		}

		if role == nil {
			c.log.Error("failed to find role for user, falling back to creating a new one",
				slog.String("target", target.ID),
				slog.String("caller", caller.ID))
		}
	}

	if role == nil {
		newRole, err := c.bot.GuildRoleCreate(c.guild.ID,
			&discordgo.RoleParams{
				Name: target.GlobalName + "'s Role",
			},
		)

		if err != nil {
			c.log.Error("failed to create role for user",
				slog.Any("error", err),
				slog.String("target", target.ID),
				slog.String("caller", caller.ID))

			return nil, c.bot.InteractionRespond(c.data.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Could not create role for user: " + target.ID,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}

		role = newRole

		c.log.Info("created role for user",
			slog.String("caller", caller.ID),
			slog.String("target", target.ID),
			slog.String("role", role.ID))

		if err := r.backend.SetRole(c.ctx, c.guild.ID, target.ID, role.ID); err != nil {
			c.log.Error("failed to set role for user",
				slog.Any("error", err),
				slog.String("target", target.ID),
				slog.String("caller", caller.ID))

			return nil, c.bot.InteractionRespond(c.data.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Could not store role for user: " + target.ID,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}

		if err := c.bot.GuildMemberRoleAdd(c.guild.ID, target.ID, role.ID); err != nil {
			c.log.Error("failed to add role to user",
				slog.Any("error", err),
				slog.String("target", target.ID),
				slog.String("caller", caller.ID))

			return nil, c.bot.InteractionRespond(c.data.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Could not add role to user: " + target.ID,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	}

	return role, nil
}

func (c *RoleUpdateContext) updatePersonalRoleName(role *discordgo.Role, name string) error {

	_, err := c.bot.GuildRoleEdit(c.guild.ID, role.ID, &discordgo.RoleParams{
		Name: name,
	})

	return err
}

func (c *RoleUpdateContext) updatePersonalRoleColor(role *discordgo.Role, input string) error {

	color, err := common.ParseTextToColorInt(input)
	if err != nil {
		return err
	}

	_, err = c.bot.GuildRoleEdit(c.guild.ID, role.ID, &discordgo.RoleParams{
		Color: &color,
	})

	return err
}
