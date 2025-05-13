package reduce

import (
	"image"
	"image/color"
)

func ReduceBits(r, g, b uint8, bits uint8) (uint8, uint8, uint8) {
    if bits > 8 {
        bits = 8 // Maksymalna liczba bitów to 8
    }

    mask := uint8(^(0xff >> bits))
	// Przyklad:
	// Zredukowana liczba ma byc np 4 bitowa, 
	// bits = 4,
	// 0xff >> 4 = 0x0f => ^0x0f = 0xf0

    reducedR := r & mask 
    reducedG := g & mask
    reducedB := b & mask

    return reducedR, reducedG, reducedB
}

func ApplyBitReductionToImage(img image.Image, bits uint8) *image.RGBA {
    bounds := img.Bounds()
    reducedImg := image.NewRGBA(bounds)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            reducedR, reducedG, reducedB := ReduceBits(uint8(r>>8), uint8(g>>8), uint8(b>>8), bits)
            reducedImg.Set(x, y, color.RGBA{R: reducedR, G: reducedG, B: reducedB, A: uint8(a >> 8)})
        }
    }
    return reducedImg
}