package grayscale

import (
	"image"
	"image/color"
)

func ConvertToGrayscale(r, g, b uint8) uint8 {
	// Formula for grayscale: 0.299*R + 0.587*G + 0.114*B
	gray := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return uint8(gray)
}

func ApplyGrayscaleToImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	grayImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := ConvertToGrayscale(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			grayColor := color.RGBA{R: gray, G: gray, B: gray, A: 255}
			grayImg.Set(x, y, grayColor)
		}
	}
	return grayImg
}