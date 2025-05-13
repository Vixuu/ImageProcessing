package hsl

import (
    "image"
    "image/color"
    "math"
)

func RGBToHSL(r, g, b uint8) (float64, float64, float64) {
    rf := float64(r) / 255.0
    gf := float64(g) / 255.0
    bf := float64(b) / 255.0

    max := math.Max(rf, math.Max(gf, bf))
    min := math.Min(rf, math.Min(gf, bf))
    delta := max - min

    // Calculate Hue
    var h float64
    if delta == 0 {
        h = 0
    } else if max == rf {
        h = math.Mod((gf-bf)/delta, 6)
    } else if max == gf {
        h = (bf-rf)/delta + 2
    } else {
        h = (rf-gf)/delta + 4
    }
    h *= 60
    if h < 0 {
        h += 360
    }

    // Calculate Lightness
    l := (max + min) / 2

    // Calculate Saturation
    var s float64
    if delta == 0 {
        s = 0
    } else {
        s = delta / (1 - math.Abs(2*l-1))
    }

    return h, s, l
}

func ConvertImageToHSL(img image.Image) [][][3]float64 {
    bounds := img.Bounds()
    width, height := bounds.Max.X, bounds.Max.Y

    hslImage := make([][][3]float64, height)
    for y := 0; y < height; y++ {
        hslImage[y] = make([][3]float64, width)
        for x := 0; x < width; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            h, s, l := RGBToHSL(uint8(r>>8), uint8(g>>8), uint8(b>>8))
            hslImage[y][x] = [3]float64{h, s, l}
        }
    }

    return hslImage
}

func HSLToRGB(h, s, l float64) (uint8, uint8, uint8) {
    c := (1 - math.Abs(2*l-1)) * s
    x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
    m := l - c/2

    var rf, gf, bf float64
    switch {
    case h >= 0 && h < 60:
        rf, gf, bf = c, x, 0
    case h >= 60 && h < 120:
        rf, gf, bf = x, c, 0
    case h >= 120 && h < 180:
        rf, gf, bf = 0, c, x
    case h >= 180 && h < 240:
        rf, gf, bf = 0, x, c
    case h >= 240 && h < 300:
        rf, gf, bf = x, 0, c
    default:
        rf, gf, bf = c, 0, x
    }

    r := uint8((rf + m) * 255)
    g := uint8((gf + m) * 255)
    b := uint8((bf + m) * 255)
    return r, g, b
}

func ConvertHSLToRGBImage(hslImage [][][3]float64) image.Image {
    height := len(hslImage)
    width := len(hslImage[0])
    rgbImg := image.NewRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            h, s, l := hslImage[y][x][0], hslImage[y][x][1], hslImage[y][x][2]
            r, g, b := HSLToRGB(h, s, l)
            rgbImg.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
        }
    }

    return rgbImg
}