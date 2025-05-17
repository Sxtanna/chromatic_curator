package imaging

import (
	"image"
	"image/color"
	"image/draw"
)

func Draw(img *image.RGBA, clr color.RGBA, rect image.Rectangle) {
	draw.Draw(img, rect, &image.Uniform{C: clr}, image.Point{}, draw.Src)
}
