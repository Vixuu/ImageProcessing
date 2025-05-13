package flip

import "image"

func FlipVertical(img image.Image) *image.RGBA {
    bounds := img.Bounds()
    flippedImg := image.NewRGBA(bounds)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            // Odbicie w pionie: zamiana wierszy
            flippedImg.Set(x, bounds.Max.Y-y-1, img.At(x, y))
        }
    }

    return flippedImg
}

func FlipHorizontal(img image.Image) *image.RGBA {
    bounds := img.Bounds()
    flippedImg := image.NewRGBA(bounds)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            // Odbicie w poziomie: zamiana kolumn
            flippedImg.Set(bounds.Max.X-x-1, y, img.At(x, y))
        }
    }

    return flippedImg
}