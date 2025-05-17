package cmds

import (
	"bytes"
	"fmt"
	"github.com/Sxtanna/chromatic_curator/internal/app/discord/data"
	"github.com/Sxtanna/chromatic_curator/internal/common"
	"github.com/Sxtanna/chromatic_curator/internal/system/imaging"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"strings"
)

// PaletteCommand represents a command to generate color palettes
type PaletteCommand struct {
	BaseCommand
}

// NewPaletteCommand creates a new color palette generation command
func NewPaletteCommand() *PaletteCommand {
	return &PaletteCommand{
		BaseCommand: BaseCommand{
			Name:        "palette",
			Description: "Generate color palettes of different types",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "type",
					Description: "The type of palette to generate",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  common.PaletteTypeMonochromatic.DisplayName(),
							Value: common.PaletteTypeMonochromatic.String(),
						},
						{
							Name:  common.PaletteTypeComplementary.DisplayName(),
							Value: common.PaletteTypeComplementary.String(),
						},
						{
							Name:  common.PaletteTypeSplitComplementary.DisplayName(),
							Value: common.PaletteTypeSplitComplementary.String(),
						},
						{
							Name:  common.PaletteTypeAnalogous.DisplayName(),
							Value: common.PaletteTypeAnalogous.String(),
						},
						{
							Name:  common.PaletteTypeTriadic.DisplayName(),
							Value: common.PaletteTypeTriadic.String(),
						},
						{
							Name:  common.PaletteTypeTetradic.DisplayName(),
							Value: common.PaletteTypeTetradic.String(),
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "color",
					Description: "The base color for the palette (name or hex code)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "count",
					Description: "Number palette layers to generate (1-10)",
					Required:    false,
					MinValue:    &[]float64{1}[0],
					MaxValue:    3,
				},
			},
		},
	}
}

// Execute handles the command execution
func (c *PaletteCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate, logger *slog.Logger) error {
	// Get the palette type from the options
	typeOption := GetOptionByName(i.Interaction, "type")
	if typeOption == nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Palette type is required",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
	paletteTypeStr := typeOption.StringValue()

	// Convert string to PaletteType enum
	paletteType, err := common.PaletteTypeFromString(paletteTypeStr)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	// Get the color name/hex from the options
	colorOption := GetOptionByName(i.Interaction, "color")
	if colorOption == nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Base color is required",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	colorName := colorOption.StringValue()
	randomColor := strings.EqualFold(colorName, "random")

	if randomColor {
		randomColorItem := common.ColorsAndNames[rand.IntN(len(common.ColorsAndNames))]
		colorName = randomColorItem.Name
	}

	colorInt, err := common.ParseTextToColorInt(colorName)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not parse color: " + colorName,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	// Get color count
	countOption := GetOptionByName(i.Interaction, "count")
	colorCount := 1 // Default

	if countOption != nil {
		colorCount = int(countOption.IntValue())
	}

	// First, acknowledge the interaction with a "thinking" response
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error("Failed to send deferred response", slog.Any("error", err))
		return err
	}

	// Generate the palette
	paletteColors, err := common.GeneratePalette(colorInt, paletteType, colorCount*5)
	if err != nil {
		logger.Error("Failed to generate palette", slog.Any("error", err))
		_, errMsg := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Failed to generate palette: " + err.Error(),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return errMsg
	}

	// Generate the palette image
	imageData, err := imaging.GenerateColorImage(colorInt, paletteColors)
	if err != nil {
		logger.Error("Failed to generate palette image", slog.Any("error", err))
		_, errMsg := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Failed to generate palette image",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return errMsg
	}

	// Create a generation object
	generation := &data.ColorGeneration{
		Input:     colorOption.StringValue(),
		ColorName: colorName,
		ColorInt:  colorInt,
		ImageData: imageData,
	}

	// Generate the embed
	generation.Embed = generatePaletteEmbed(generation, paletteColors, paletteType)

	// Generate a UUID for the image and store it in the cache
	imageID := data.SaveGeneration(generation)

	// Create a custom ID for the share button that includes the UUID
	shareButtonID := fmt.Sprintf("share_color:%s", imageID)

	// Send a follow-up message with the image and a share button
	tempMessage, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Flags:  discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{generation.Embed},
		Files: []*discordgo.File{
			{
				Name:   "palette_preview.png",
				Reader: bytes.NewReader(imageData),
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Share to Channel",
						Style:    discordgo.PrimaryButton,
						CustomID: shareButtonID,
					},
				},
			},
		},
	})

	if err != nil {
		logger.Error("Failed to send follow-up message", slog.Any("error", err))
		return err
	}

	generation.TempMsgID = tempMessage.ID

	return nil
}

// generatePaletteEmbed creates an embed for the palette
func generatePaletteEmbed(generation *data.ColorGeneration, paletteColors []common.ColorDistance, paletteType common.PaletteType) *discordgo.MessageEmbed {
	// Convert the int color to RGB
	r, g, b := common.IntToRGB(generation.ColorInt)

	// Create the embed for the palette
	embed := &discordgo.MessageEmbed{
		Title:       "Color Palette",
		Description: fmt.Sprintf("**%s** palette based on **%s**", paletteType.DisplayName(), generation.ColorName),
		Color:       generation.ColorInt,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Base Color",
				Value:  fmt.Sprintf("`%s` - `#%02X%02X%02X` - `RGB(%d, %d, %d)`", generation.ColorName, r, g, b, r, g, b),
				Inline: false,
			},
		},
	}

	// Add palette colors info to the embed
	if len(paletteColors) > 0 {
		minNameLength := 0

		for _, similar := range paletteColors {
			minNameLength = max(minNameLength, len(similar.Name))
		}

		minNameLength += 2

		var paletteColorsText strings.Builder

		// Create a table header
		paletteColorsText.WriteString("```\n")
		paletteColorsText.WriteString(fmt.Sprintf("   %-"+strconv.Itoa(minNameLength)+"s %-10s %-10s\n", "Color Name", "Hex Code", "RGB"))
		paletteColorsText.WriteString(fmt.Sprintf("   %-"+strconv.Itoa(minNameLength)+"s %-10s %-10s\n", "----------", "--------", "-------------"))

		// Add each color as a row in the table
		for i, color := range paletteColors {
			// Convert the color int to RGB
			r, g, b := common.IntToRGB(color.ColorInt)

			// Add row to table with index number
			paletteColorsText.WriteString(fmt.Sprintf("%-2d %-"+strconv.Itoa(minNameLength)+"s %-10s (%d, %d, %d)\n",
				i+1, color.Name, fmt.Sprintf("#%02X%02X%02X", r, g, b), r, g, b))
		}
		paletteColorsText.WriteString("```")

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Palette Colors (%d) (* = est)", len(paletteColors)),
			Value:  paletteColorsText.String(),
			Inline: false,
		})
	}

	return embed
}
