# Image Compression

This project provides tools and utilities for image compression, helping to reduce file sizes while maintaining acceptable image quality.

## Features

- Multiple compression algorithms support
- Batch processing capabilities
- Quality adjustment options
- Support for common image formats (PNG, JPEG, etc.)

## Installation

```bash
# Clone the repository
git clone https://github.com/uzzalcse/image_compression.git

# Navigate to project directory
cd image_compression

# Install dependencies
go mod tidy

# Run the project 
go run main.go

```



## Usage

### Basic example:

To compress images you have to keep images to the folder named `in_images`. Keeping images in this folder you have to run the project. 
After running the project your out put images will be save in the  `out_images` folder. In the projects root directory.

## Requirements

- Go version 1.23.2
- pkg-config
- libvips-dev


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.