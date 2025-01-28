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
		log.Printf("Failed to read input directory: %v", err)
		return
	}

	// Process each file in the directory
	for _, file := range files {
		if !file.IsDir() && isImage(file.Name()) {
			// Get the full path to the image
			imagePath := filepath.Join(inputDir, file.Name())

			fmt.Printf("\nProcessing image: %s\n", imagePath)
			
			// Print original image dimensions
			if err := printImageDimensions(imagePath); err != nil {
				log.Printf("Error getting dimensions for %s: %v", imagePath, err)
				// Continue to next operations even if this fails
			}

			// Run each compression method independently
			// If one fails, continue with others
			benchmarkCompression("disintegration/imaging", compressWithImaging, imagePath)
			benchmarkCompression("h2non/bimg", compressWithBimg, imagePath)
			benchmarkCompression("nfnt/resize", compressWithResize, imagePath)
			benchmarkCompression("rwcarlsen/goexif/exif", compressWithExif, imagePath)
			benchmarkCompression("anthonynsimon/bild", compressWithBild, imagePath)
		}
	}
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func printImageDimensions(imagePath string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	fmt.Printf("Image Dimensions for %s: %d x %d\n", imagePath, img.Bounds().Dx(), img.Bounds().Dy())
	return nil
}

func benchmarkCompression(methodName string, compressFunc func(string) (string, error), imagePath string) {
	fmt.Printf("\nAttempting compression with %s\n", methodName)
	
	startTime := time.Now()

	// Perform compression and handle potential error
	compressedFile, err := compressFunc(imagePath)
	if err != nil {
		log.Printf("Error in %s compression: %v", methodName, err)
		return // Skip benchmarking if compression failed
	}

	duration := time.Since(startTime)

	// Get file sizes
	originalFile, err := os.Stat(imagePath)
	if err != nil {
		log.Printf("Error getting original file info: %v", err)
		return
	}

	compressedFileInfo, err := os.Stat(compressedFile)
	if err != nil {
		log.Printf("Error getting compressed file info: %v", err)
		return
	}

	fmt.Printf("[%s Benchmark]\n", methodName)
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Original file size: %d bytes\n", originalFile.Size())
	fmt.Printf("Compressed file size: %d bytes\n", compressedFileInfo.Size())
}

func compressWithImaging(imagePath string) (string, error) {
	fmt.Println("[Using disintegration/imaging]")

	src, err := imaging.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}

	outputFile := generateOutputFileName(imagePath, "disintegration_imaging")

	if err := saveImageToFolder(src, outputFile); err != nil {
		return "", fmt.Errorf("failed to save compressed image: %v", err)
	}

	return outputFile, nil
}

func compressWithBimg(imagePath string) (string, error) {
	fmt.Println("[Using h2non/bimg]")

	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read input image: %v", err)
	}

	image := bimg.NewImage(imageData)

	options := bimg.Options{
		Quality: 60,
		Type:    bimg.JPEG,
	}

	compressedImage, err := image.Process(options)
	if err != nil {
		return "", fmt.Errorf("failed to process image: %v", err)
	}

	outputFile := generateOutputFileName(imagePath, "h2non_bimg")

	if err := saveImageToFolder(compressedImage, outputFile); err != nil {
		return "", fmt.Errorf("failed to save compressed image: %v", err)
	}

	return outputFile, nil
}

func compressWithResize(imagePath string) (string, error) {
	fmt.Println("[Using nfnt/resize]")

	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	width := uint(img.Bounds().Dx())
	height := uint(img.Bounds().Dy())
	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)

	outputFile := generateOutputFileName(imagePath, "nfnt_resize")

	if err := saveImageToFolder(resizedImg, outputFile); err != nil {
		return "", fmt.Errorf("failed to save resized image: %v", err)
	}

	return outputFile, nil
}

func compressWithExif(imagePath string) (string, error) {
	fmt.Println("[Using rwcarlsen/goexif/exif]")

	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	exifData, err := exif.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode EXIF data: %v", err)
	}

	cameraModel, _ := exifData.Get(exif.Model)
	fmt.Printf("Camera Model: %v\n", cameraModel)

	file.Seek(0, 0)
	img, err := imaging.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image for compression: %v", err)
	}

	outputFile := generateOutputFileName(imagePath, "rwcarlsen_goexif_exif")

	if err := saveImageToFolder(img, outputFile); err != nil {
		return "", fmt.Errorf("failed to save compressed image: %v", err)
	}

	return outputFile, nil
}

func compressWithBild(imagePath string) (string, error) {
	fmt.Println("[Using anthonynsimon/bild]")

	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	width := float64(img.Bounds().Dx())
	height := float64(img.Bounds().Dy())

	resized := transform.Resize(img, int(width), int(height), transform.Linear)

	outputFile := generateOutputFileName(imagePath, "anthonynsimon_bild")

	if err := saveImageToFolder(resized, outputFile); err != nil {
		return "", fmt.Errorf("failed to save compressed image: %v", err)
	}

	return outputFile, nil
}

func generateOutputFileName(imagePath, packageName string) string {
	baseName := strings.TrimPrefix(imagePath, "in_images/")
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
	return fmt.Sprintf("out_images/%s_compressed_%s%s", baseName, packageName, filepath.Ext(imagePath))
}

func saveImageToFolder(img interface{}, outputFile string) error {
	if err := os.MkdirAll("out_images", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create 'out_images' folder: %v", err)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	switch v := img.(type) {
	case image.Image:
		if err := jpeg.Encode(outFile, v, &jpeg.Options{Quality: 60}); err != nil {
			return fmt.Errorf("failed to encode image: %v", err)
		}
	case []byte:
		if _, err := outFile.Write(v); err != nil {
			return fmt.Errorf("failed to write compressed image bytes: %v", err)
		}
	default:
		return fmt.Errorf("unsupported image type: %T", v)
	}

	return nil
}