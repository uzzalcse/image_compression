package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	_"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"gopkg.in/h2non/bimg.v1"
	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	// Specify the input directory
	inputDir := "in_images"

	// Get all image files in the 'in_images' directory
	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	// Process each file in the directory
	for _, file := range files {
		if !file.IsDir() && isImage(file.Name()) {
			// Get the full path to the image
			imagePath := filepath.Join(inputDir, file.Name())

			// Print original image dimensions
			printImageDimensions(imagePath)

			// Benchmark Imaging Compression
			benchmarkCompression("disintegration/imaging", compressWithImaging, imagePath)

			// Benchmark Bimg Compression
			benchmarkCompression("h2non/bimg", compressWithBimg, imagePath)

			// Benchmark nfnt/resize Compression
			benchmarkCompression("nfnt/resize", compressWithResize, imagePath)

			// Benchmark EXIF Compression
			benchmarkCompression("rwcarlsen/goexif/exif", compressWithExif, imagePath)
			// Benchmark Bild Compression
			benchmarkCompression("anthonynsimon/bild", compressWithBild, imagePath)
		}
	}
}

// Function to check if a file is an image
func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
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

	// Generate dynamic output file name (without in_images prefix)
	outputFile := generateOutputFileName(imagePath, "disintegration_imaging")

	// Save the compressed image with reduced quality to the out_images folder
	err = saveImageToFolder(src, outputFile)
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

	// Generate dynamic output file name (without in_images prefix)
	outputFile := generateOutputFileName(imagePath, "h2non_bimg")

	// Save the compressed image to the out_images folder
	err = saveImageToFolder(compressedImage, outputFile)
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

	// Generate dynamic output file name (without in_images prefix)
	outputFile := generateOutputFileName(imagePath, "nfnt_resize")

	// Save the resized image to the out_images folder
	err = saveImageToFolder(resizedImg, outputFile)
	if err != nil {
		log.Fatalf("Failed to save resized image: %v", err)
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

	// Generate dynamic output file name (without in_images prefix)
	outputFile := generateOutputFileName(imagePath, "rwcarlsen_goexif_exif")

	// Save the compressed image to the out_images folder
	err = saveImageToFolder(img, outputFile)
	if err != nil {
		log.Fatalf("Failed to save compressed image: %v", err)
	}

	fmt.Printf("Compressed image with EXIF metadata saved as: %s\n", outputFile)
	return outputFile
}

// New function for compression using "anthonynsimon/bild"

func compressWithBild(imagePath string) string {
	fmt.Println("\n[Using anthonynsimon/bild]")

	// Load the image using standard image package first for better performance
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	// Calculate new dimensions (90% of original)
	width := float64(img.Bounds().Dx()) * 0.9
	height := float64(img.Bounds().Dy()) * 0.9

	// Use bild's transform with ResizeStrategyFit for better performance
	resized := transform.Resize(img, int(width), int(height), transform.Linear)

	// Generate output filename
	outputFile := generateOutputFileName(imagePath, "anthonynsimon_bild")

	// Create output file
	out, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer out.Close()

	// Use jpeg.Encode directly with quality option
	err = jpeg.Encode(out, resized, &jpeg.Options{Quality: 70})
	if err != nil {
		log.Fatalf("Failed to encode compressed image: %v", err)
	}

	fmt.Printf("Compressed image saved as: %s\n", outputFile)
	return outputFile
}

// Generate a dynamic output file name based on the original file name and package name
func generateOutputFileName(imagePath, packageName string) string {
	// Extract the base name (without extension) of the original image
	baseName := strings.TrimPrefix(imagePath, "in_images/")  // Remove "in_images/" prefix
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName)) // Remove file extension
	// Return the formatted output file name
	return fmt.Sprintf("out_images/%s_compressed_%s%s", baseName, packageName, filepath.Ext(imagePath))
}

// Save the image to the specified folder
func saveImageToFolder(img interface{}, outputFile string) error {
	// Ensure the 'out_images' folder exists
	err := os.MkdirAll("out_images", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create 'out_images' folder: %v", err)
	}

	// Create the output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Save the image based on its type (image.Image for imaging and resize, []byte for bimg)
	switch v := img.(type) {
	case image.Image:
		// Encode and save the image if it's of type image.Image
		err = jpeg.Encode(outFile, v, &jpeg.Options{Quality: 70})
		if err != nil {
			return fmt.Errorf("failed to encode image: %v", err)
		}
	case []byte:
		// Save the image if it's of type []byte (for bimg)
		_, err = outFile.Write(v)
		if err != nil {
			return fmt.Errorf("failed to write compressed image bytes: %v", err)
		}
	default:
		return fmt.Errorf("unsupported image type: %T", v)
	}

	return nil
}
