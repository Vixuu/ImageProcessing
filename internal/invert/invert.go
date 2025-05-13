package invert

import (
	"image"
	"image/color"
)

func InvertColor(r, g, b uint8) (uint8, uint8, uint8) {
	return 255 - r, 255 - g, 255 - b
}

func ApplyColorInversionToImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	invertedImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			invertedR, invertedG, invertedB := InvertColor(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			invertedColor := color.RGBA{R: invertedR, G: invertedG, B: invertedB, A: 255}
			invertedImg.Set(x, y, invertedColor)
		}
	}
	return invertedImg
}