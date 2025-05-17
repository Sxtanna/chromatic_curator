package cmds

import (
	"bytes"
	"fmt"
	"github.com/Sxtanna/chromatic_curator/internal/common"
	"github.com/Sxtanna/chromatic_curator/internal/system/imaging"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
)

type ColorGeneration struct {
	Input     string
	ColorName string
	ColorInt  int
	ImageData []byte
	Embed     *discordgo.MessageEmbed
	TempMsgID string
}

// Global cache for storing color images
var (
	ColorImageCache      = make(map[string]*ColorGeneration)
	ColorImageCacheMutex sync.RWMutex
)

// StoreColorImage stores an image in the cache with the given ID
func StoreColorImage(gen *ColorGeneration) string {
	// Generate a UUID for the image
	id := uuid.New().String()

	// Store the image in the cache
	ColorImageCacheMutex.Lock()
	ColorImageCache[id] = gen
	ColorImageCacheMutex.Unlock()

	return id
}

// ColorCommand represents a command to preview colors
type ColorCommand struct {
	BaseCommand
}

// NewColorCommand creates a new color preview command
func NewColorCommand() *ColorCommand {
	return &ColorCommand{
		BaseCommand: BaseCommand{
			Name:        "color",
			Description: "Preview colors from the available color set",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "The name or hex code of the color to preview",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "range",
					Description: "Number of similar colors to show (1-10)",
					Required:    false,
					MinValue:    &[]float64{1}[0],
					MaxValue:    10,
				},
			},
		},
	}
}

// Execute handles the command execution
func (c *ColorCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate, logger *slog.Logger) error {
	// Get the color name/hex from the options
	colorOption := GetOptionByName(i.Interaction, "name")
	if colorOption == nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Color name or hex code is required",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	colorName := colorOption.StringValue()
	randomColors := strings.EqualFold(colorName, "random")

	if randomColors {
		randomColor := common.ColorsAndNames[rand.IntN(len(common.ColorsAndNames))]
		colorName = randomColor.Name
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

	generation := &ColorGeneration{
		Input:     colorOption.StringValue(),
		ColorName: colorName,
		ColorInt:  colorInt,
	}

	// Check if range option is provided
	rangeOption := GetOptionByName(i.Interaction, "range")
	var rangeValue = 0
	if rangeOption != nil {
		rangeValue = int(rangeOption.IntValue())
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

	// Find similar colors if range is specified
	var similarColors []common.ColorDistance
	if rangeValue > 0 {
		if !randomColors {
			similarColors = common.FindSimilarColors(colorInt, rangeValue)
		} else {
			similarColors = make([]common.ColorDistance, rangeValue)

			for i := 0; i < rangeValue; i++ {
				randomColor := common.ColorsAndNames[rand.IntN(len(common.ColorsAndNames))]

				itemColorInt, err := common.ParseTextToColorInt(randomColor.Color)
				if err != nil || itemColorInt == colorInt {
					i--
					continue
				}

				similarColors[i] = common.ColorDistance{
					Name:     randomColor.Name,
					ColorInt: itemColorInt,
					Distance: 0,
				}
			}
		}
	}

	// Generate the color preview image
	imageData, err := imaging.GenerateColorImage(colorInt, similarColors)
	if err != nil {
		logger.Error("Failed to generate color image", slog.Any("error", err))
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to generate color image",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	generation.ImageData = imageData
	generation.Embed = generateColorGenerationEmbed(generation, similarColors, randomColors)

	// Generate a UUID for the image and store it in the cache
	imageID := StoreColorImage(generation)

	// Create a custom ID for the share button that includes the UUID
	shareButtonID := fmt.Sprintf("share_color:%s", imageID)

	// Then, send a follow-up message with the image and a share button
	tempMessage, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Flags:  discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{generation.Embed},
		Files: []*discordgo.File{
			{
				Name:   "color_preview.png",
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

func generateColorGenerationEmbed(generation *ColorGeneration, similarColors []common.ColorDistance, randomColors bool) *discordgo.MessageEmbed {
	// Convert the int color to RGB
	r, g, b := common.IntToRGB(generation.ColorInt)

	// Create the embed for the main color
	embed := &discordgo.MessageEmbed{
		Title:       "Color Preview",
		Description: fmt.Sprintf("Preview of color: **%s**", generation.ColorName),
		Color:       generation.ColorInt,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Hex Code",
				Value:  "`" + fmt.Sprintf("#%02X%02X%02X", r, g, b) + "`",
				Inline: true,
			},
			{
				Name:   "RGB",
				Value:  fmt.Sprintf("`(%d, %d, %d)`", r, g, b),
				Inline: true,
			},
		},
	}

	// Add similar colors info to the embed if applicable
	if len(similarColors) > 0 {

		minNameLength := 0

		for _, similar := range similarColors {
			minNameLength = max(minNameLength, len(similar.Name))
		}

		minNameLength += 2

		var similarColorsText strings.Builder

		// Create a table header
		similarColorsText.WriteString("```\n")
		similarColorsText.WriteString(fmt.Sprintf("   %-"+strconv.Itoa(minNameLength)+"s %-10s %-10s\n", "Color Name", "Hex Code", "RGB"))
		similarColorsText.WriteString(fmt.Sprintf("   %-"+strconv.Itoa(minNameLength)+"s %-10s %-10s\n", "----------", "--------", "-------------"))

		// Add each color as a row in the table
		for i, similar := range similarColors {
			// Convert the color int to RGB
			r, g, b := common.IntToRGB(similar.ColorInt)

			// Add row to table with index number
			similarColorsText.WriteString(fmt.Sprintf("%-2d %-"+strconv.Itoa(minNameLength)+"s %-10s (%d, %d, %d)\n",
				i+1, similar.Name, fmt.Sprintf("#%02X%02X%02X", r, g, b), r, g, b))
		}
		similarColorsText.WriteString("```")

		embedTitle := "Similar Colors"
		if randomColors {
			embedTitle = "Random Colors"
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s (%d)", embedTitle, len(similarColors)),
			Value:  similarColorsText.String(),
			Inline: false,
		})
	}

	return embed
}
