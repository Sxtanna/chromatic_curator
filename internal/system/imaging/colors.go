package imaging

import (
	"bytes"
	"github.com/Sxtanna/chromatic_curator/internal/common"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

// Define image dimensions
const (
	MAIN_COLOR_SIZE     = 150 // Main color square is larger
	ALTS_COLOR_SIZE     = 100 // Size for similar colors
	SQUARE_PADDING      = 10  // Space between elements
	SQUARE_BORDER_SIZE  = 2   // Size of the border around color squares
	MAX_SQUARES_PER_ROW = 5   // Maximum number of colors per row
)

// GenerateColorImage creates an image showing the main color and similar colors
func GenerateColorImage(mainColor int, similarColors []common.ColorDistance) ([]byte, error) {
	rowCount := 1 // At least one row for the main color

	if len(similarColors) > 0 {
		rowCount += (len(similarColors) + MAX_SQUARES_PER_ROW - 1) / MAX_SQUARES_PER_ROW
	}

	// Calculate width and height
	fullImageWidth := max(MAIN_COLOR_SIZE+2*SQUARE_PADDING, min(len(similarColors), MAX_SQUARES_PER_ROW)*ALTS_COLOR_SIZE+(min(len(similarColors), MAX_SQUARES_PER_ROW)+1)*SQUARE_PADDING)
	fullImageHeight := SQUARE_PADDING + MAIN_COLOR_SIZE + SQUARE_PADDING // Space for main color
	if len(similarColors) > 0 {
		fullImageHeight += (rowCount - 1) * (ALTS_COLOR_SIZE + SQUARE_PADDING) // Space for similar colors
	}

	// Create a new RGBA image
	imageData := image.NewRGBA(image.Rect(0, 0, fullImageWidth, fullImageHeight))

	draw.Draw(imageData, imageData.Bounds(), &image.Uniform{C: color.RGBA{R: 0, G: 0, B: 0, A: 0}}, image.Point{}, draw.Src)

	// Draw the main color square first
	mainR, mainG, mainB := common.IntToRGB(mainColor)
	darkR, darkG, darkB := common.SlightlyDarker(mainR, mainG, mainB)

	DrawSquareWithBorder(imageData,
		color.RGBA{R: mainR, G: mainG, B: mainB, A: 255},
		color.RGBA{R: darkR, G: darkG, B: darkB, A: 255},
		(fullImageWidth-MAIN_COLOR_SIZE)/2,
		SQUARE_PADDING,
		MAIN_COLOR_SIZE,
		SQUARE_BORDER_SIZE)

	// Draw similar colors
	if len(similarColors) > 0 {
		startY := SQUARE_PADDING + MAIN_COLOR_SIZE + SQUARE_PADDING

		for i, colorDist := range similarColors {
			// Get color components
			distR, distG, distB := common.IntToRGB(colorDist.ColorInt)
			darkR, darkG, darkB := common.SlightlyDarker(distR, distG, distB)

			x := SQUARE_PADDING + (i%MAX_SQUARES_PER_ROW)*(ALTS_COLOR_SIZE+SQUARE_PADDING)
			y := startY + (i/MAX_SQUARES_PER_ROW)*(ALTS_COLOR_SIZE+SQUARE_PADDING)

			DrawSquareWithBorder(imageData,
				color.RGBA{R: distR, G: distG, B: distB, A: 255},
				color.RGBA{R: darkR, G: darkG, B: darkB, A: 255},
				x,
				y,
				ALTS_COLOR_SIZE,
				SQUARE_BORDER_SIZE)

			// Choose contrasting color for the text (white or black depending on color brightness)
			textColor := color.RGBA{R: 255, G: 255, B: 255, A: 255} // Default to white

			// Use black for light colors (simple brightness calculation)
			brightness := (int(distR) + int(distG) + int(distB)) / 3
			if brightness > 128 {
				textColor = color.RGBA{A: 255}
			}

			// Draw the index number directly on the color square
			DrawNumber(imageData, i+1, x+ALTS_COLOR_SIZE/2, y+ALTS_COLOR_SIZE/2, textColor)
		}
	}

	// Encode the image to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, imageData); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
