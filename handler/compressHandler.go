package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cepot-blip/opencv-go-project/utils"
)

type CompressRequest struct {
	InputPath string `json:"input_path"`
	Quality   int    `json:"quality"`
}

func CompressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// DECODE JSON REQUEST
	var req CompressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	// OPEN IMAGE FILE
	file, err := os.Open(req.InputPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open image file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// GET THE BASE NAME OF THE INPUT FIlE
	inputFileName := filepath.Base(req.InputPath)

	// SPECIPY OUTPUT PATH INSADE "ASSETS" FOLDER WITH JPEG FORMAT
	outputPath := filepath.Join("assets", fmt.Sprintf("compressed_%s", inputFileName))

	// COMPRESS IMAGE
	err = utils.CompressImage(file, outputPath, req.Quality)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error compressing image: %v", err), http.StatusInternalServerError)
		return
	}

	// OPEN AND SERVE THE COMPRESSED IMAGE FILE
	compressedFile, err := os.Open(outputPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open compressed image file: %v", err), http.StatusInternalServerError)
		return
	}
	defer compressedFile.Close()

	// SET CONTENT TYPE
	w.Header().Set("Content-Type", "image/jpeg")

	// COPY THE COMPRESSED IMAGE FILE TO THE RESPONSE WRITER
	_, err = io.Copy(w, compressedFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to copy compressed image: %v", err), http.StatusInternalServerError)
		return
	}
}
