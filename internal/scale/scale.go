package scale

import (
	"image"
	"image/color"
	"math"
)


func ScaleCoords(originalWidth, originalHeight int, scaleX, scaleY float64, x, y int) (int, int) {
    // Przelicz współrzędne
    newX := int(float64(x) * scaleX)
    newY := int(float64(y) * scaleY)

    return newX, newY
}

func ReverseScaleCoords(scaledWidth, scaledHeight int, scaleX, scaleY float64, x, y int) (int, int) {
    // Przelicz współrzędne odwrotne
    originalX := int(math.Round(float64(x) / scaleX))
    originalY := int(math.Round(float64(y) / scaleY))

    return originalX, originalY
}

func ResizeImage(img image.Image, newWidth, newHeight int) *image.RGBA {
    bounds := img.Bounds()
    originalWidth := bounds.Dx()
    originalHeight := bounds.Dy()

    scaleX := float64(originalWidth) / float64(newWidth)
    scaleY := float64(originalHeight) / float64(newHeight)

    resizedImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

    for y := 0; y < newHeight; y++ {
        for x := 0; x < newWidth; x++ {
            newX, newY := ScaleCoords(originalWidth, originalHeight, scaleX, scaleY, x, y)

            if newX >= 0 && newX < originalWidth && newY >= 0 && newY < originalHeight {
                origColor := img.At(newX, newY)
                resizedImg.Set(x, y, origColor)
            } else {
                resizedImg.Set(x, y, color.RGBA{0, 0, 0, 255})
            }
        }
    }

    return resizedImg
}