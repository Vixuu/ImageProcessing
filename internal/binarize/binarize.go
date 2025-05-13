package binarize

import (
	"image"
	"image/color"
	"zadanie/1/internal/grayscale"
)

func BinarizeColor(r, g, b uint8, threshold uint8) uint8 {
    gray := grayscale.ConvertToGrayscale(r, g, b)
    if gray >= threshold {
        return 1 // Bright pixel
    }
    return 0 // Dark pixel
}

func ApplyBinarizationToImage(img image.Image, threshold uint8) *image.Gray {
    bounds := img.Bounds()
    binaryImg := image.NewGray(bounds)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            binaryValue := BinarizeColor(uint8(r>>8), uint8(g>>8), uint8(b>>8), threshold)
            if binaryValue == 1 {
                binaryImg.SetGray(x, y, color.Gray{Y: 255}) // White
            } else {
                binaryImg.SetGray(x, y, color.Gray{Y: 0}) // Black
            }
        }
    }
    return binaryImg
}