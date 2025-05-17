package common

import (
	"emperror.dev/errors"
	"math"
)

// GeneratePalette generates a color palette based on the specified type
func GeneratePalette(baseColor int, paletteType PaletteType, colorCount int) ([]ColorDistance, error) {
	// Convert the base color to HSV for easier manipulation
	r, g, b := IntToRGB(baseColor)
	h, s, v := RGBToHSV(r, g, b)

	// Generate the palette based on the type
	var palette []ColorDistance

	switch paletteType {
	case PaletteTypeMonochromatic:
		palette = GenerateMonochromaticPalette(h, s, v, colorCount)
	case PaletteTypeComplementary:
		palette = GenerateComplementaryPalette(h, s, v, colorCount)
	case PaletteTypeSplitComplementary:
		palette = GenerateSplitComplementaryPalette(h, s, v, colorCount)
	case PaletteTypeAnalogous:
		palette = GenerateAnalogousPalette(h, s, v, colorCount)
	case PaletteTypeTriadic:
		palette = GenerateTriadicPalette(h, s, v, colorCount)
	case PaletteTypeTetradic:
		palette = GenerateTetradicPalette(h, s, v, colorCount)
	default:
		return nil, errors.Errorf("unknown palette type: %s", paletteType)
	}

	return palette, nil
}

// GenerateMonochromaticPalette generates a monochromatic palette
func GenerateMonochromaticPalette(h, s, v float64, numColors int) []ColorDistance {
	palette := make([]ColorDistance, numColors)

	// For monochromatic, we vary the saturation and value
	for i := 0; i < numColors; i++ {
		// Vary saturation and value based on position
		newS := math.Max(0.1, math.Min(1.0, s*(0.5+float64(i)/float64(numColors))))
		newV := math.Max(0.2, math.Min(1.0, v*(0.6+float64(numColors-i)/float64(numColors))))

		r, g, b := HSVToRGB(h, newS, newV)
		palette[i] = FindExactOrClosestNamedColor(r, g, b)
	}

	return palette
}

// GenerateComplementaryPalette generates a complementary palette
func GenerateComplementaryPalette(h, s, v float64, numColors int) []ColorDistance {
	palette := make([]ColorDistance, numColors)

	// Complementary color is 180 degrees away
	complementaryH := math.Mod(h+180, 360)

	// Distribute colors between the base and its complement
	for i := 0; i < numColors; i++ {
		var newH, newS, newV float64

		if i < numColors/2 {
			// Colors closer to the base
			newH = h
			newS = math.Max(0.1, math.Min(1.0, s*(0.7+float64(i)/float64(numColors))))
			newV = math.Max(0.2, math.Min(1.0, v*(0.7+float64(i)/float64(numColors))))
		} else {
			// Colors closer to the complement
			newH = complementaryH
			newS = math.Max(0.1, math.Min(1.0, s*(0.7+float64(numColors-i)/float64(numColors))))
			newV = math.Max(0.2, math.Min(1.0, v*(0.7+float64(i)/float64(numColors))))
		}

		r, g, b := HSVToRGB(newH, newS, newV)
		palette[i] = FindExactOrClosestNamedColor(r, g, b)
	}

	return palette
}

// GenerateSplitComplementaryPalette generates a split complementary palette
func GenerateSplitComplementaryPalette(h, s, v float64, numColors int) []ColorDistance {
	palette := make([]ColorDistance, numColors)

	// Split complementary colors are 150 and 210 degrees away
	splitComplement1 := math.Mod(h+150, 360)
	splitComplement2 := math.Mod(h+210, 360)

	// Distribute colors among the three main hues
	for i := 0; i < numColors; i++ {
		var newH float64

		if i < numColors/3 {
			newH = h
		} else if i < 2*numColors/3 {
			newH = splitComplement1
		} else {
			newH = splitComplement2
		}

		// Vary saturation and value slightly for variety
		newS := math.Max(0.1, math.Min(1.0, s*(0.7+float64(i%3)/3.0)))
		newV := math.Max(0.2, math.Min(1.0, v*(0.7+float64(i%3)/3.0)))

		r, g, b := HSVToRGB(newH, newS, newV)
		palette[i] = FindExactOrClosestNamedColor(r, g, b)
	}

	return palette
}

// GenerateAnalogousPalette generates an analogous palette
func GenerateAnalogousPalette(h, s, v float64, numColors int) []ColorDistance {
	palette := make([]ColorDistance, numColors)

	// Analogous colors are within 30 degrees on either side
	angleRange := 60.0 // 30 degrees on each side

	for i := 0; i < numColors; i++ {
		// Distribute hues evenly across the range
		newH := math.Mod(h-angleRange/2+angleRange*float64(i)/float64(numColors-1), 360)

		// Vary saturation and value slightly for variety
		newS := math.Max(0.1, math.Min(1.0, s*(0.7+float64(i%3)/3.0)))
		newV := math.Max(0.2, math.Min(1.0, v*(0.7+float64(i%3)/3.0)))

		r, g, b := HSVToRGB(newH, newS, newV)
		palette[i] = FindExactOrClosestNamedColor(r, g, b)
	}

	return palette
}

// GenerateTriadicPalette generates a triadic palette
func GenerateTriadicPalette(h, s, v float64, numColors int) []ColorDistance {
	palette := make([]ColorDistance, numColors)

	// Triadic colors are 120 degrees apart
	triad1 := math.Mod(h+120, 360)
	triad2 := math.Mod(h+240, 360)

	// Distribute colors among the three main hues
	for i := 0; i < numColors; i++ {
		var newH float64

		if i < numColors/3 {
			newH = h
		} else if i < 2*numColors/3 {
			newH = triad1
		} else {
			newH = triad2
		}

		// Vary saturation and value slightly for variety
		newS := math.Max(0.1, math.Min(1.0, s*(0.7+float64(i%3)/3.0)))
		newV := math.Max(0.2, math.Min(1.0, v*(0.7+float64(i%3)/3.0)))

		r, g, b := HSVToRGB(newH, newS, newV)
		palette[i] = FindExactOrClosestNamedColor(r, g, b)
	}

	return palette
}

// GenerateTetradicPalette generates a tetradic palette
func GenerateTetradicPalette(h, s, v float64, numColors int) []ColorDistance {
	palette := make([]ColorDistance, numColors)

	// Tetradic colors are 90 degrees apart
	tetrad1 := math.Mod(h+90, 360)
	tetrad2 := math.Mod(h+180, 360)
	tetrad3 := math.Mod(h+270, 360)

	// Distribute colors among the four main hues
	for i := 0; i < numColors; i++ {
		var newH float64

		if i < numColors/4 {
			newH = h
		} else if i < numColors/2 {
			newH = tetrad1
		} else if i < 3*numColors/4 {
			newH = tetrad2
		} else {
			newH = tetrad3
		}

		// Vary saturation and value slightly for variety
		newS := math.Max(0.1, math.Min(1.0, s*(0.7+float64(i%4)/4.0)))
		newV := math.Max(0.2, math.Min(1.0, v*(0.7+float64(i%4)/4.0)))

		r, g, b := HSVToRGB(newH, newS, newV)
		palette[i] = FindExactOrClosestNamedColor(r, g, b)
	}

	return palette
}
