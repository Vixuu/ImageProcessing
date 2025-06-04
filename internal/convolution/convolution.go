package convolution

import (
	"image"
	"image/color"
	"math"
)

type PaddingType string

const (
    None      PaddingType = "none"
    Zero      PaddingType = "zero"
    Replicate PaddingType = "replicate"
)

// ConvertToGrayMatrix konwertuje obraz na macierz float64 (szarość)
func ConvertToGrayMatrix(img image.Image) [][]float64 {
    bounds := img.Bounds()
    width, height := bounds.Max.X, bounds.Max.Y
    data := make([][]float64, height)
    for y := 0; y < height; y++ {
        data[y] = make([]float64, width)
        for x := 0; x < width; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            gray := float64((r + g + b) / 3 >> 8)
            data[y][x] = gray
        }
    }
    return data
}

// ConvertGrayMatrixToImage konwertuje macierz float64 na obraz Gray
func ConvertGrayMatrixToImage(data [][]float64) *image.Gray {
    height := len(data)
    width := len(data[0])
    out := image.NewGray(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            val := uint8(math.Max(0, math.Min(255, data[y][x])))
            out.SetGray(x, y, color.Gray{Y: val})
        }
    }
    return out
}

// normalizeKernel normalizuje kernel
func normalizeKernel(kernel [][]float64) [][]float64 {
    sum := 0.0
    for _, row := range kernel {
        for _, val := range row {
            sum += val
        }
    }
    if sum == 0 {
        return kernel
    }
    out := make([][]float64, len(kernel))
    for i := range kernel {
        out[i] = make([]float64, len(kernel[i]))
        for j := range kernel[i] {
            out[i][j] = kernel[i][j] / sum
        }
    }
    return out
}

// getPixel pobiera piksel z paddingiem
func getPixel(img [][]float64, x, y int, padding PaddingType) float64 {
    h, w := len(img), len(img[0])
    if x >= 0 && x < h && y >= 0 && y < w {
        return img[x][y]
    }
    switch padding {
    case Zero:
        return 0
    case Replicate:
        ix := int(math.Max(0, math.Min(float64(x), float64(h-1))))
        iy := int(math.Max(0, math.Min(float64(y), float64(w-1))))
        return img[ix][iy]
    default:
        return math.NaN()
    }
}

// Convolve wykonuje konwolucję na macierzy obrazu
func Convolve(img [][]float64, kernel [][]float64, padding PaddingType) [][]float64 {
    kernel = normalizeKernel(kernel)
    kH, kW := len(kernel), len(kernel[0])
    centerY, centerX := kH/2, kW/2
    h, w := len(img), len(img[0])
    out := make([][]float64, h)
    for y := 0; y < h; y++ {
        out[y] = make([]float64, w)
        for x := 0; x < w; x++ {
            sum := 0.0
            valid := true
            for i := 0; i < kH; i++ {
                for j := 0; j < kW; j++ {
                    yy := y + i - centerY
                    xx := x + j - centerX
                    val := getPixel(img, yy, xx, padding)
                    if math.IsNaN(val) {
                        valid = false
                        break
                    }
                    sum += val * kernel[i][j]
                }
            }
            if padding == None && !valid {
                out[y][x] = img[y][x]
            } else {
                out[y][x] = sum
            }
        }
    }
    return out
}