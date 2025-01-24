// package main

// import (
// 	"fmt"
// 	"log"
// 	_"os"
// 	"path/filepath"

// 	"github.com/disintegration/imaging"
// )

// func main() {
// 	// Specify the original image file
// 	originalFile := "image_1.jpeg" // Replace with your actual image file name

// 	// Open the original image
// 	src, err := imaging.Open(originalFile)
// 	if err != nil {
// 		log.Fatalf("Failed to open image: %v", err)
// 	}

// 	// Set the output file name with "reduced_" prefix
// 	dir, file := filepath.Split(originalFile)
// 	ext := filepath.Ext(file)
// 	name := file[:len(file)-len(ext)]
// 	compressedFile := filepath.Join(dir, fmt.Sprintf("reduced_%s%s", name, ext))

// 	// Save the compressed image with a reduced quality (e.g., 85 for JPEG)
// 	err = imaging.Save(src, compressedFile, imaging.JPEGQuality(85))
// 	if err != nil {
// 		log.Fatalf("Failed to save compressed image: %v", err)
// 	}

// 	fmt.Printf("Compressed image saved as: %s\n", compressedFile)
// }


package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func main() {
	// Specify the original image file
	originalFile := "image_1.jpeg" // Replace with your actual image file name

	// Open the original image
	src, err := imaging.Open(originalFile)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}

	// Get the dimensions of the original image
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	fmt.Printf("Original image dimensions: %dx%d\n", width, height)

	// Set the output file name with "reduced_" prefix
	dir, file := filepath.Split(originalFile)
	ext := filepath.Ext(file)
	name := file[:len(file)-len(ext)]
	compressedFile := filepath.Join(dir, fmt.Sprintf("reduced_%s%s", name, ext))

	// Save the compressed image with the same dimensions and reduced quality
	err = imaging.Save(src, compressedFile, imaging.JPEGQuality(85))
	if err != nil {
		log.Fatalf("Failed to save compressed image: %v", err)
	}

	fmt.Printf("Compressed image saved as: %s\n", compressedFile)

	// Verify the dimensions of the compressed image
	reducedImage, err := imaging.Open(compressedFile)
	if err != nil {
		log.Fatalf("Failed to open compressed image: %v", err)
	}

	reducedWidth := reducedImage.Bounds().Dx()
	reducedHeight := reducedImage.Bounds().Dy()

	fmt.Printf("Compressed image dimensions: %dx%d\n", reducedWidth, reducedHeight)

	if width == reducedWidth && height == reducedHeight {
		fmt.Println("The dimensions of the original and compressed images are the same.")
	} else {
		fmt.Println("The dimensions of the compressed image differ from the original.")
	}
}
