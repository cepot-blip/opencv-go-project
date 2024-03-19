package utils

import (
	"fmt"
	"image"
	"os"

	"gocv.io/x/gocv"
)

// RESIZE IMAGE RESIZES THE IMAGE USING GOCV
func ResizeImage(InputPath *os.File, outputPath string, width, height int) error {
	img := gocv.IMRead(InputPath.Name(), gocv.IMReadColor)
	if img.Empty() {
		return fmt.Errorf("failed to decode image")
	}
	defer img.Close()

	resized := gocv.NewMat()
	defer resized.Close()

	gocv.Resize(img, &resized, image.Point{width, height}, 0, 0, gocv.InterpolationDefault)

	// WRITES AN IMAGE WITHOUT INCLUDING AN OPTION FOR JPEG QUALITY
	if !gocv.IMWrite(outputPath, resized) {
		return fmt.Errorf("failed to write image")
	}

	return nil
}
