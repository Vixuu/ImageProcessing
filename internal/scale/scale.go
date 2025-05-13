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

func ResizeImageBilinear(src image.Image, newW, newH int) image.Image {
    dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
    srcBounds := src.Bounds()
    srcW := srcBounds.Dx()
    srcH := srcBounds.Dy()

    scaleX := float64(srcW) / float64(newW)
    scaleY := float64(srcH) / float64(newH)

    for y := 0; y < newH; y++ {
        for x := 0; x < newW; x++ {
            // Pozycja w oryginalnym obrazie
            fx := float64(x)*scaleX + 0.5*(scaleX-1)
            fy := float64(y)*scaleY + 0.5*(scaleY-1)

            x0 := int(math.Floor(fx))
            y0 := int(math.Floor(fy))
            x1 := x0 + 1
            y1 := y0 + 1

            wx := fx - float64(x0)
            wy := fy - float64(y0)

            // Granice obrazu
            x0 = clamp(x0, 0, srcW-1)
            x1 = clamp(x1, 0, srcW-1)
            y0 = clamp(y0, 0, srcH-1)
            y1 = clamp(y1, 0, srcH-1)

            c00 := src.At(srcBounds.Min.X+x0, srcBounds.Min.Y+y0)
            c10 := src.At(srcBounds.Min.X+x1, srcBounds.Min.Y+y0)
            c01 := src.At(srcBounds.Min.X+x0, srcBounds.Min.Y+y1)
            c11 := src.At(srcBounds.Min.X+x1, srcBounds.Min.Y+y1)

            r := bilinear(
                float64(getR(c00)), float64(getR(c10)),
                float64(getR(c01)), float64(getR(c11)),
                wx, wy,
            )
            g := bilinear(
                float64(getG(c00)), float64(getG(c10)),
                float64(getG(c01)), float64(getG(c11)),
                wx, wy,
            )
            b := bilinear(
                float64(getB(c00)), float64(getB(c10)),
                float64(getB(c01)), float64(getB(c11)),
                wx, wy,
            )
            a := bilinear(
                float64(getA(c00)), float64(getA(c10)),
                float64(getA(c01)), float64(getA(c11)),
                wx, wy,
            )

            dst.Set(x, y, color.NRGBA{
                R: uint8(clamp(int(r+0.5), 0, 255)),
                G: uint8(clamp(int(g+0.5), 0, 255)),
                B: uint8(clamp(int(b+0.5), 0, 255)),
                A: uint8(clamp(int(a+0.5), 0, 255)),
            })
        }
    }
    return dst
}

func bilinear(c00, c10, c01, c11, wx, wy float64) float64 {
    return (1-wx)*(1-wy)*c00 +
        wx*(1-wy)*c10 +
        (1-wx)*wy*c01 +
        wx*wy*c11
}

func clamp(x, min, max int) int {
    if x < min {
        return min
    }
    if x > max {
        return max
    }
    return x
}

func getR(c color.Color) uint32 {
    r, _, _, _ := c.RGBA()
    return r >> 8
}
func getG(c color.Color) uint32 {
    _, g, _, _ := c.RGBA()
    return g >> 8
}
func getB(c color.Color) uint32 {
    _, _, b, _ := c.RGBA()
    return b >> 8
}
func getA(c color.Color) uint32 {
    _, _, _, a := c.RGBA()
    return a >> 8
}