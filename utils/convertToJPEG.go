package utils

import (
	"fmt"
	"os"

	"gocv.io/x/gocv"
)

// CONVERT TO JPEG CONVERTS IMAGES FROM PNG TO JPEG FORMAT USING GOCV
func ConvertToJPEG(input *os.File, outputPath string) error {
	img := gocv.IMRead(input.Name(), gocv.IMReadColor)
	if img.Empty() {
		return fmt.Errorf("failed to decode image")
	}
	defer img.Close()

	// COPY THE IMAGE TO SET ITS JPEG PROPERTIES
	jpegImage := gocv.NewMat()
	defer jpegImage.Close()
	img.CopyTo(&jpegImage)

	// WRITE IMAGES IN JPEG FORMAT WITHOUT ADDITIONAL OPTIONS
	if !gocv.IMWrite(outputPath, jpegImage) {
		return fmt.Errorf("failed to write image as JPEG")
	}

	return nil
}
