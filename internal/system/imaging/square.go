package imaging

import (
	"image"
	"image/color"
)

func DrawSquareWithBorder(img *image.RGBA, squareColor, squareBorderColor color.RGBA, xCoord, yCoord int, squareSize, borderSize int) {
	// Draw border (slightly darker than the color)
	Draw(img, squareBorderColor, image.Rect(
		xCoord-borderSize,
		yCoord-borderSize,
		xCoord+squareSize+borderSize,
		yCoord+squareSize+borderSize,
	))

	// Draw main color
	Draw(img, squareColor, image.Rect(
		xCoord,
		yCoord,
		xCoord+squareSize,
		yCoord+squareSize,
	))
}
