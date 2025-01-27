package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"gopkg.in/h2non/bimg.v1"
	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	// Specify the original image file
	originalFile := "img_2.jpeg" // Replace with your actual image file name

	// Print original image dimensions
	printImageDimensions(originalFile)

	// Benchmark Imaging Compression
	benchmarkCompression("disintegration/imaging", compressWithImaging, originalFile)

	// Benchmark Bimg Compression
	benchmarkCompression("h2non/bimg", compressWithBimg, originalFile)

	// Benchmark nfnt/resize Compression
	benchmarkCompression("nfnt/resize", compressWithResize, originalFile)

	// Benchmark EXIF Compression
	benchmarkCompression("rwcarlsen/goexif/exif", compressWithExif, originalFile)
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

// Benchmark function to measure time, file size, and quality
func benchmarkCompression(methodName string, compressFunc func(string) string, imagePath string) {
	startTime := time.Now()

	// Perform compression
	compressedFile := compressFunc(imagePath)

	// Measure time taken for compression
	duration := time.Since(startTime)

	// Get the file size of the original image
	originalFile, err := os.Stat(imagePath)
	if err != nil {
		log.Fatalf("Failed to get original image file info: %v", err)
	}

	// Get the file size of the compressed image
	compressedFileInfo, err := os.Stat(compressedFile)
	if err != nil {
		log.Fatalf("Failed to get compressed image file info: %v", err)
	}

	// Display the benchmarking results
	fmt.Printf("\n[%s Benchmark]\n", methodName)
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Original file size: %d bytes\n", originalFile.Size())
	fmt.Printf("Compressed file size: %d bytes\n", compressedFileInfo.Size())
}

// Compression using "disintegration/imaging"
func compressWithImaging(imagePath string) string {
	fmt.Println("\n[Using disintegration/imaging]")

	// Open the original image
	src, err := imaging.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}

	// Generate dynamic output file name
	outputFile := generateOutputFileName(imagePath, "disintegration_imaging")

	// Save the compressed image with reduced quality
	err = imaging.Save(src, outputFile, imaging.JPEGQuality(75))
	if err != nil {
		log.Fatalf("Failed to save compressed image: %v", err)
	}

	fmt.Printf("Compressed image saved as: %s\n", outputFile)
	return outputFile
}

// Compression using "h2non/bimg"
func compressWithBimg(imagePath string) string {
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

	// Generate dynamic output file name
	outputFile := generateOutputFileName(imagePath, "h2non_bimg")

	// Save the compressed image
	err = os.WriteFile(outputFile, compressedImage, 0644)
	if err != nil {
		log.Fatalf("Failed to save compressed image: %v", err)
	}

	fmt.Printf("Compressed image saved as: %s\n", outputFile)
	return outputFile
}

// Compression using "nfnt/resize"
func compressWithResize(imagePath string) string {
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

	// Generate dynamic output file name
	outputFile := generateOutputFileName(imagePath, "nfnt_resize")

	// Save the resized image
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
	return outputFile
}

// Compression with EXIF Metadata Handling using "rwcarlsen/goexif/exif"
func compressWithExif(imagePath string) string {
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

	// Generate dynamic output file name
	outputFile := generateOutputFileName(imagePath, "rwcarlsen_goexif_exif")

	// Save the compressed image
	err = imaging.Save(img, outputFile, imaging.JPEGQuality(70))
	if err != nil {
		log.Fatalf("Failed to save compressed image: %v", err)
	}

	fmt.Printf("Compressed image with EXIF metadata saved as: %s\n", outputFile)
	return outputFile
}

// Generate a dynamic output file name based on the original file name and package name
func generateOutputFileName(imagePath, packageName string) string {
	// Extract the base name (without extension) of the original image
	baseName := strings.TrimSuffix(imagePath, ".jpeg")
	// Return the formatted output file name
	return fmt.Sprintf("%s_compressed_%s.jpeg", baseName, packageName)
}
