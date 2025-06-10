package morphology

import (
    "image"
    "image/color"
)

// Zamienia obraz na macierz binarną (0 lub 1)
func ImageToBinaryMatrix(img image.Image, threshold uint8) [][]uint8 {
    bounds := img.Bounds()
    w, h := bounds.Dx(), bounds.Dy()
    mat := make([][]uint8, h)
    for y := 0; y < h; y++ {
        mat[y] = make([]uint8, w)
        for x := 0; x < w; x++ {
            r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
            gray := uint8((r>>8 + g>>8 + b>>8) / 3)
            if gray > threshold {
                mat[y][x] = 1
            } else {
                mat[y][x] = 0
            }
        }
    }
    return mat
}

// Zamienia macierz binarną na obraz
func BinaryMatrixToImage(mat [][]uint8) *image.Gray {
    h := len(mat)
    w := len(mat[0])
    img := image.NewGray(image.Rect(0, 0, w, h))
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            val := uint8(0)
            if mat[y][x] > 0 {
                val = 255
            }
            img.SetGray(x, y, color.Gray{Y: val})
        }
    }
    return img
}

// Erozja
func Erode(bin [][]uint8, kernel [][]int) [][]uint8 {
    h, w := len(bin), len(bin[0])
    kh, kw := len(kernel), len(kernel[0])
    cy, cx := kh/2, kw/2
    out := make([][]uint8, h)
    for y := 0; y < h; y++ {
        out[y] = make([]uint8, w)
        for x := 0; x < w; x++ {
            match := true
            for ky := 0; ky < kh; ky++ {
                for kx := 0; kx < kw; kx++ {
                    iy, ix := y+ky-cy, x+kx-cx
                    if iy < 0 || iy >= h || ix < 0 || ix >= w {
                        match = false
                        break
                    }
                    if kernel[ky][kx] == 1 && bin[iy][ix] == 0 {
                        match = false
                        break
                    }
                }
            }
            if match {
                out[y][x] = 1
            }
        }
    }
    return out
}

// Dylatacja
func Dilate(bin [][]uint8, kernel [][]int) [][]uint8 {
    h, w := len(bin), len(bin[0])
    kh, kw := len(kernel), len(kernel[0])
    cy, cx := kh/2, kw/2
    out := make([][]uint8, h)
    for y := 0; y < h; y++ {
        out[y] = make([]uint8, w)
        for x := 0; x < w; x++ {
            found := false
            for ky := 0; ky < kh; ky++ {
                for kx := 0; kx < kw; kx++ {
                    iy, ix := y+ky-cy, x+kx-cx
                    if iy < 0 || iy >= h || ix < 0 || ix >= w {
                        continue
                    }
                    if kernel[ky][kx] == 1 && bin[iy][ix] == 1 {
                        found = true
                        break
                    }
                }
                if found {
                    break
                }
            }
            if found {
                out[y][x] = 1
            }
        }
    }
    return out
}

// Otwarcie: erozja, potem dylatacja
func Open(bin [][]uint8, kernel [][]int) [][]uint8 {
    return Dilate(Erode(bin, kernel), kernel)
}

// Zamknięcie: dylatacja, potem erozja
func Close(bin [][]uint8, kernel [][]int) [][]uint8 {
    return Erode(Dilate(bin, kernel), kernel)
}

// Hit-or-miss transformacja
func HitOrMiss(bin [][]uint8, hitKernel, missKernel [][]int) [][]uint8 {
    h, w := len(bin), len(bin[0])
    kh, kw := len(hitKernel), len(hitKernel[0])
    cy, cx := kh/2, kw/2
    out := make([][]uint8, h)
    for y := 0; y < h; y++ {
        out[y] = make([]uint8, w)
        for x := 0; x < w; x++ {
            hit := true
            for ky := 0; ky < kh; ky++ {
                for kx := 0; kx < kw; kx++ {
                    iy, ix := y+ky-cy, x+kx-cx
                    if iy < 0 || iy >= h || ix < 0 || ix >= w {
                        hit = false
                        break
                    }
                    if hitKernel[ky][kx] == 1 && bin[iy][ix] != 1 {
                        hit = false
                        break
                    }
                    if missKernel[ky][kx] == 1 && bin[iy][ix] != 0 {
                        hit = false
                        break
                    }
                }
            }
            if hit {
                out[y][x] = 1
            }
        }
    }
    return out
}

// Szkieletyzacja (prosta iteracyjna, aż do wyzerowania)
func Skeletonize(bin [][]uint8, kernel [][]int) [][]uint8 {
    prev := make([][]uint8, len(bin))
    for i := range bin {
        prev[i] = make([]uint8, len(bin[0]))
        copy(prev[i], bin[i])
    }
    tmp := make([][]uint8, len(bin))
    for i := range bin {
        tmp[i] = make([]uint8, len(bin[0]))
    }
    for {
        eroded := Erode(prev, kernel)
        hitmiss := HitOrMiss(eroded, kernel, invertKernel(kernel))
        for y := range prev {
            for x := range prev[0] {
                tmp[y][x] = eroded[y][x]
                if hitmiss[y][x] == 1 {
                    tmp[y][x] = 0
                }
            }
        }
        if equal(tmp, prev) {
            break
        }
        for y := range prev {
            copy(prev[y], tmp[y])
        }
    }
    return prev
}

// Pomocnicza: odwraca kernel (1->0, 0->1)
func invertKernel(kernel [][]int) [][]int {
    h := len(kernel)
    w := len(kernel[0])
    out := make([][]int, h)
    for y := 0; y < h; y++ {
        out[y] = make([]int, w)
        for x := 0; x < w; x++ {
            if kernel[y][x] == 1 {
                out[y][x] = 0
            } else {
                out[y][x] = 1
            }
        }
    }
    return out
}

// Pomocnicza: porównuje dwie macierze
func equal(a, b [][]uint8) bool {
    for y := range a {
        for x := range a[0] {
            if a[y][x] != b[y][x] {
                return false
            }
        }
    }
    return true
}