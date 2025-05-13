package histogram

import (
    "image"
    "image/color"
    _ "image/jpeg"
    _ "image/png"
    "math"
    "os"
    "path/filepath"

    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
)

// GenerateHistogram generuje i zapisuje znormalizowany histogram jasności lub RGB.
// path - ścieżka do pliku obrazu
// tryb - "jasnosc" lub "rgb"
// outputDir - folder, do którego zostanie zapisany histogram
func GenerateHistogram(path string, tryb string, outputDir string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return err
    }

    switch tryb {
    case "jasnosc":
        hist := make([]float64, 256)
        var max float64
        bounds := img.Bounds()
        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
            for x := bounds.Min.X; x < bounds.Max.X; x++ {
                r, g, b, _ := img.At(x, y).RGBA()
                rr := float64(r >> 8)
                gg := float64(g >> 8)
                bb := float64(b >> 8)
                brightness := 0.299*rr + 0.587*gg + 0.114*bb
                idx := int(math.Round(brightness))
                if idx > 255 {
                    idx = 255
                }
                hist[idx]++
                if hist[idx] > max {
                    max = hist[idx]
                }
            }
        }
        for i := range hist {
            hist[i] /= max
        }
        p := plot.New()
        p.Title.Text = "Histogram jasności"
        p.X.Label.Text = "Jasność"
        p.Y.Label.Text = "Znormalizowana liczba pikseli"
        pts := make(plotter.XYs, 256)
        for i := 0; i < 256; i++ {
            pts[i].X = float64(i)
            pts[i].Y = hist[i]
        }
        line, err := plotter.NewLine(pts)
        if err != nil {
            return err
        }
        line.Color = color.Black
        p.Add(line)
        filename := filepath.Join(outputDir, "histogram_jasnosc.png")
        return p.Save(8*vg.Inch, 4*vg.Inch, filename)

    case "rgb":
        histR := make([]float64, 256)
        histG := make([]float64, 256)
        histB := make([]float64, 256)
        var maxR, maxG, maxB float64
        bounds := img.Bounds()
        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
            for x := bounds.Min.X; x < bounds.Max.X; x++ {
                r, g, b, _ := img.At(x, y).RGBA()
                ri := int(r >> 8)
                gi := int(g >> 8)
                bi := int(b >> 8)
                histR[ri]++
                histG[gi]++
                histB[bi]++
                if histR[ri] > maxR {
                    maxR = histR[ri]
                }
                if histG[gi] > maxG {
                    maxG = histG[gi]
                }
                if histB[bi] > maxB {
                    maxB = histB[bi]
                }
            }
        }
        for i := 0; i < 256; i++ {
            histR[i] /= maxR
            histG[i] /= maxG
            histB[i] /= maxB
        }
        p := plot.New()
        p.Title.Text = "Histogram RGB"
        p.X.Label.Text = "Wartość kanału"
        p.Y.Label.Text = "Znormalizowana liczba pikseli"
        ptsR := make(plotter.XYs, 256)
        ptsG := make(plotter.XYs, 256)
        ptsB := make(plotter.XYs, 256)
        for i := 0; i < 256; i++ {
            ptsR[i].X = float64(i)
            ptsR[i].Y = histR[i]
            ptsG[i].X = float64(i)
            ptsG[i].Y = histG[i]
            ptsB[i].X = float64(i)
            ptsB[i].Y = histB[i]
        }
        lineR, _ := plotter.NewLine(ptsR)
        lineG, _ := plotter.NewLine(ptsG)
        lineB, _ := plotter.NewLine(ptsB)
        lineR.Color = color.RGBA{255, 0, 0, 255}
        lineG.Color = color.RGBA{0, 255, 0, 255}
        lineB.Color = color.RGBA{0, 0, 255, 255}
        p.Add(lineR, lineG, lineB)
        filename := filepath.Join(outputDir, "histogram_rgb.png")
        return p.Save(8*vg.Inch, 4*vg.Inch, filename)
    default:
        return nil
    }
}