package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	_"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"gopkg.in/h2non/bimg.v1"
	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	// Specify the original image file
	originalFile := "image_1.jpeg" // Replace with your actual image file name

	// Print original image dimensions
	printImageDimensions(originalFile)

	// Imaging Compression
	compressWithImaging(originalFile)

	// Bimg Compression
	compressWithBimg(originalFile)

	// nfnt/resize Compression
	compressWithResize(originalFile)

	// EXIF Metadata Handling
	compressWithExif(originalFile)
}

// Function to print the dimensions of an image
func printImageDimensions(imagePath string) {
	// Open the image
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()

	// Decode the image to get dimensions
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	// Print the image dimensions
	fmt.Printf("Image Dimensions for %s: %d x %d\n", imagePath, img.Bounds().Dx(), img.Bounds().Dy())
}

// Compression using "disintegration/imaging"
func compressWithImaging(imagePath string) {
	fmt.Println("\n[Using disintegration/imaging]")

	// Open the original image
	src, err := imaging.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}

	// Save the compressed image with reduced quality
	outputFile := "compressed_imaging.jpeg"
	err = imaging.Save(src, outputFile, imaging.JPEGQuality(75))
	if err != nil {
		log.Fatalf("Failed to save compressed image: %v", err)
	}


	fmt.Printf("Compressed image saved as: %s\n", outputFile)
}

// Compression using "h2non/bimg"
func compressWithBimg(imagePath string) {
	fmt.Println("\n[Using h2non/bimg]")

	// Read the input image
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("Failed to read input image: %v", err)
	}

	// Create a new bimg instance
	image := bimg.NewImage(imageData)

	// Compress the image
	options := bimg.Options{
		Quality: 70, // Reduce quality
		Type:    bimg.JPEG,
	}

	compressedImage, err := image.Process(options)
	if err != nil {
		log.Fatalf("Failed to process image: %v", err)
	}

	// Save the compressed image
	outputFile := "compressed_bimg.jpeg"
	err = os.WriteFile(outputFile, compressedImage, 0644)
	if err != nil {
		log.Fatalf("Failed to save compressed image: %v", err)
	}


	fmt.Printf("Compressed image saved as: %s\n", outputFile)
}

// Compression using "nfnt/resize"
func compressWithResize(imagePath string) {
	fmt.Println("\n[Using nfnt/resize]")

	// Open the original image
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	// Resize the image to half its original dimensions
	width := uint(img.Bounds().Dx())
	height := uint(img.Bounds().Dy())
	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)

	// Save the resized image
	outputFile := "compressed_resize.jpeg"
	outFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 70})
	if err != nil {
		log.Fatalf("Failed to encode resized image: %v", err)
	}


	fmt.Printf("Resized and compressed image saved as: %s\n", outputFile)
}

// Compression with EXIF Metadata Handling using "rwcarlsen/goexif/exif"
func compressWithExif(imagePath string) {
	fmt.Println("\n[Using rwcarlsen/goexif/exif]")

	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()

	// Extract EXIF metadata
	exifData, err := exif.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode EXIF data: %v", err)
	}

	// Get camera model from metadata
	cameraModel, _ := exifData.Get(exif.Model)
	fmt.Printf("Camera Model: %v\n", cameraModel)

	// Re-open the image for compression (since EXIF decoding reads the file pointer)
	file.Seek(0, 0)
	img, err := imaging.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image for compression: %v", err)
	}

	// Save the compressed image
	outputFile := "compressed_exif.jpeg"
	err = imaging.Save(img, outputFile, imaging.JPEGQuality(70))
	if err != nil {
		log.Fatalf("Failed to save compressed image: %v", err)
	}


	fmt.Printf("Compressed image with EXIF metadata saved as: %s\n", outputFile)
}
