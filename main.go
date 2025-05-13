package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"
	"zadanie/1/internal/binarize"
	"zadanie/1/internal/flip"
	"zadanie/1/internal/grayscale"
	"zadanie/1/internal/histogram"
	"zadanie/1/internal/hsl"
	"zadanie/1/internal/invert"
	"zadanie/1/internal/reduce"
	"zadanie/1/internal/rotate"
	"zadanie/1/internal/scale"
)

// LoadImage loads an image from a file
func LoadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// SaveImage saves an image to a file
func SaveImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetPixelRGB prints the RGB values of a pixel at the given position
func GetPixelRGB(img image.Image, x, y int) {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		fmt.Println("Pixel out of bounds")
		return
	}

	r, g, b, _ := img.At(x, y).RGBA()
	fmt.Printf("Pixel at (%d, %d): R=%d, G=%d, B=%d\n", x, y, r>>8, g>>8, b>>8)
}

func main() {

	start := time.Now()

	os.MkdirAll("output/hist", os.ModePerm)


	// Wczytanie obrazu
	img, err := LoadImage("input.jpg")
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	// Wartosci RGB piksela
	x, y := 50, 50
	GetPixelRGB(img, x, y)

	//Grayscale
	grayImg := grayscale.ApplyGrayscaleToImage(img)
	err = SaveImage(grayImg, "output/grayscale.jpg")
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	//Binarize
	binaryImg := binarize.ApplyBinarizationToImage(img, 127) // Próg jasności: 128
	err = SaveImage(binaryImg, "output/binary.jpg")
	if err != nil {
		fmt.Println("Error saving binary image:", err)
		return
	}

	// Invert
	invertedImg := invert.ApplyColorInversionToImage(img)
	err = SaveImage(invertedImg, "output/inverted.jpg")
	if err != nil {
		fmt.Println("Error saving inverted image:", err)
		return
	}

	// r, g, b := reduce.ReduceBits(255, 16, 32, 4)
	// println(r, g, b)

	// Reduce bits
	reducedImg := reduce.ApplyBitReductionToImage(img, 4)
	err = SaveImage(reducedImg, "output/reduced.jpg")
	if err != nil {
		fmt.Println("Error saving reduced image:", err)
		return
	}

	// Przekonwertuj RGB na HSL
	hslImg := hsl.ConvertImageToHSL(img)

	//  Przekonwertuj HSL na RGB
	rgbImg := hsl.ConvertHSLToRGBImage(hslImg)
	err = SaveImage(rgbImg, "output/rgb.jpg")
	if err != nil {
		fmt.Println("Error saving RGB image:", err)
		return
	}
	fmt.Println("Image processing completed successfully.")

	// Rotate
	rotatedImg := rotate.RotateImage(img, 2) // Rotate 3 times (270 degrees)
	err = SaveImage(rotatedImg, "output/rotated.jpg")
	if err != nil {
		fmt.Println("Error saving rotated image:", err)
	}

	// Flip vertically
	flippedImg := flip.FlipVertical(img)
	err = SaveImage(flippedImg, "output/flipped_vertical.jpg")
	if err != nil {
		fmt.Println("Error saving flipped image:", err)
	}

	flippedImg = flip.FlipHorizontal(img)
	err = SaveImage(flippedImg, "output/flipped_horizontal.jpg")
	if err != nil {
		fmt.Println("Error saving flipped image:", err)
	}

	// Resize the image
	resizedImg := scale.ResizeImage(img, 1000, 800)

	// Save the resized image
	err = SaveImage(resizedImg, "output/resized.jpg")
	if err != nil {
		fmt.Println("Error saving resized image:", err)
	}

	resized := scale.ResizeImageBilinear(img, 800, 600)
	err = SaveImage(resized, "output/resized_bilinear.jpg")
	if err != nil {
		fmt.Println("Error saving resized bilinear image:", err)
	}

	err = histogram.GenerateHistogram("input.jpg", "jasnosc", "output/hist")
	if err != nil {
		log.Fatal(err)
	}
	err = histogram.GenerateHistogram("input.jpg", "rgb", "output/hist")
    if err != nil {
        log.Fatal(err)
    }

	end := time.Since(start)
	fmt.Printf("Processing time: %s\n", end.Round(time.Millisecond))
}
