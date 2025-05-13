package rotate

import (
	"image"
)

func RotateImage(img image.Image, rotations int) image.Image {
    rotations = rotations % 4
    if rotations < 0 {
        rotations += 4
    }

    bounds := img.Bounds()
    var rotatedImg *image.RGBA

    switch rotations {
    case 1: // 90 degrees
        rotatedImg = image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
            for x := bounds.Min.X; x < bounds.Max.X; x++ {
                rotatedImg.Set(bounds.Max.Y-y-1, x, img.At(x, y))
            }
        }
    case 2: // 180 degrees
        rotatedImg = image.NewRGBA(bounds)
        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
            for x := bounds.Min.X; x < bounds.Max.X; x++ {
                rotatedImg.Set(bounds.Max.X-x-1, bounds.Max.Y-y-1, img.At(x, y))
            }
        }
    case 3: // 270 degrees
        rotatedImg = image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
            for x := bounds.Min.X; x < bounds.Max.X; x++ {
                rotatedImg.Set(y, bounds.Max.X-x-1, img.At(x, y))
            }
        }
    case 0: // 0 degrees 
        return img
    }

    return rotatedImg
}

