package imaging

import (
	"fmt"
	"image"
	"image/color"
)

const (
	DEFAULT_DIGIT_WIDTH  = 20
	DEFAULT_DIGIT_HEIGHT = 30
	DEFAULT_DIGIT_STROKE = 4
)

// DrawNumber draws a number on the image at the specified position
// It breaks down the number into individual digits and draws each one
func DrawNumber(img *image.RGBA, number int, xCoord, yCoord int, clr color.RGBA) {
	// Handle negative numbers
	if number < 0 {
		number = -number // Just use absolute value for now
	}

	// Convert number to string to get individual digits
	numStr := fmt.Sprintf("%d", number)

	// Calculate total width to center the number
	totalWidth := len(numStr) * DEFAULT_DIGIT_WIDTH

	// Starting x position (centered)
	startX := xCoord - totalWidth/2 + DEFAULT_DIGIT_WIDTH/2

	// Draw each digit
	for i, digitRune := range numStr {
		digit := int(digitRune - '0') // Convert rune to int
		digitX := startX + i*DEFAULT_DIGIT_WIDTH
		DrawDigit(img, digit, digitX, yCoord, clr)
	}
}

// DrawDigit draws a single digit (0-9) on the image at the specified position
func DrawDigit(img *image.RGBA, digit int, xCoord, yCoord int, clr color.RGBA) {
	// Define the size of the digit
	var (
		width  = 20
		height = 30
		stroke = 4
	)

	// Calculate the top-left corner of the digit
	xCoord = xCoord - width/2
	yCoord = yCoord - height/2

	// Draw different shapes based on the digit
	switch digit {
	case 0:
		// Left vertical line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+stroke,
			yCoord+height,
		))

		// Right vertical line
		Draw(img, clr, image.Rect(
			xCoord+width-stroke,
			yCoord,
			xCoord+width,
			yCoord+height,
		))

		// Top horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+width,
			yCoord+stroke,
		))

		// Bottom horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height-stroke,
			xCoord+width,
			yCoord+height,
		))
	case 1:
		// Draw a vertical line for 1
		Draw(img, clr, image.Rect(
			xCoord+width/2-stroke/2,
			yCoord,
			xCoord+width/2+stroke/2,
			yCoord+height,
		))
	case 2:
		// Draw top horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+width,
			yCoord+stroke,
		))

		// Draw middle horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height/2-stroke/2,
			xCoord+width,
			yCoord+height/2+stroke/2,
		))

		// Draw bottom horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height-stroke,
			xCoord+width,
			yCoord+height,
		))

		// Draw top-right vertical line
		Draw(img, clr, image.Rect(
			xCoord+width-stroke,
			yCoord,
			xCoord+width,
			yCoord+height/2,
		))

		// Draw bottom-left vertical line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height/2,
			xCoord+stroke,
			yCoord+height,
		))

	case 3:
		// Draw top horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+width,
			yCoord+stroke,
		))

		// Draw middle horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height/2-stroke/2,
			xCoord+width,
			yCoord+height/2+stroke/2,
		))

		// Draw bottom horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height-stroke,
			xCoord+width,
			yCoord+height,
		))

		// Draw right vertical line
		Draw(img, clr, image.Rect(
			xCoord+width-stroke,
			yCoord,
			xCoord+width,
			yCoord+height,
		))

	case 4:
		// Draw left vertical line (top half)
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+stroke,
			yCoord+height/2,
		))

		// Draw middle horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height/2-stroke/2,
			xCoord+width,
			yCoord+height/2+stroke/2,
		))

		// Draw right vertical line
		Draw(img, clr, image.Rect(
			xCoord+width-stroke,
			yCoord,
			xCoord+width,
			yCoord+height,
		))

	case 5:
		// Draw top horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+width,
			yCoord+stroke,
		))

		// Draw middle horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height/2-stroke/2,
			xCoord+width,
			yCoord+height/2+stroke/2,
		))

		// Draw bottom horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height-stroke,
			xCoord+width,
			yCoord+height,
		))

		// Draw top-left vertical line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+stroke,
			yCoord+height/2,
		))

		// Draw bottom-right vertical line
		Draw(img, clr, image.Rect(
			xCoord+width-stroke,
			yCoord+height/2,
			xCoord+width,
			yCoord+height,
		))

	case 6:
		// Draw left vertical line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+stroke,
			yCoord+height,
		))

		// Draw top horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+width,
			yCoord+stroke,
		))

		// Draw middle horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height/2-stroke/2,
			xCoord+width,
			yCoord+height/2+stroke/2,
		))

		// Draw bottom horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height-stroke,
			xCoord+width,
			yCoord+height,
		))

		// Draw bottom-right vertical line
		Draw(img, clr, image.Rect(
			xCoord+width-stroke,
			yCoord+height/2,
			xCoord+width,
			yCoord+height,
		))

	case 7:
		// Draw top horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+width,
			yCoord+stroke,
		))

		// Draw right vertical line
		Draw(img, clr, image.Rect(
			xCoord+width-stroke,
			yCoord,
			xCoord+width,
			yCoord+height,
		))

	case 8:
		// Draw left vertical line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+stroke,
			yCoord+height,
		))

		// Draw right vertical line
		Draw(img, clr, image.Rect(
			xCoord+width-stroke,
			yCoord,
			xCoord+width,
			yCoord+height,
		))

		// Draw top horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+width,
			yCoord+stroke,
		))

		// Draw middle horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height/2-stroke/2,
			xCoord+width,
			yCoord+height/2+stroke/2,
		))

		// Draw bottom horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height-stroke,
			xCoord+width,
			yCoord+height,
		))

	case 9:
		// Draw right vertical line
		Draw(img, clr, image.Rect(
			xCoord+width-stroke,
			yCoord,
			xCoord+width,
			yCoord+height,
		))

		// Draw top horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+width,
			yCoord+stroke,
		))

		// Draw middle horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height/2-stroke/2,
			xCoord+width,
			yCoord+height/2+stroke/2,
		))

		// Draw bottom horizontal line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord+height-stroke,
			xCoord+width,
			yCoord+height,
		))

		// Draw top-left vertical line
		Draw(img, clr, image.Rect(
			xCoord,
			yCoord,
			xCoord+stroke,
			yCoord+height/2,
		))

	// case 10 is removed as we now have DrawNumber function to handle multi-digit numbers

	default:
		// For invalid digits (not 0-9), do nothing or draw a simple indicator
		// This should not happen if DrawNumber is used correctly
		Draw(img, clr, image.Rect(
			xCoord+width/4,
			yCoord+height/4,
			xCoord+width*3/4,
			yCoord+height*3/4,
		))
	}
}
