package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// COMPRESS IMAGE PERFORMS IMAGE COMPRESSION USING FFMPEG WITH QUALITY OPTIONS
func CompressImage(input *os.File, outputPath string, quality int) error {
	cmd := exec.Command("ffmpeg", "-i", input.Name(), "-vf", fmt.Sprintf("scale=640:480"), "-q:v", fmt.Sprintf("%d", quality), outputPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error compressing image: %v", err)
	}
	return nil
}
