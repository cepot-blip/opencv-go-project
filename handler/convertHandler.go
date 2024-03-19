package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cepot-blip/opencv-go-project/utils"
)

type ConvertRequest struct {
	PNGPath  string `json:"png_path"`
	JPEGPath string `json:"jpeg_path"`
}

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// DECODE JSON REQUEST
	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	// OPEN PNG FILE
	pngFile, err := os.Open(req.PNGPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open PNG file: %v", err), http.StatusInternalServerError)
		return
	}
	defer pngFile.Close()

	// SPECIPY OUTPUT PATH
	jpegOutputPath := req.JPEGPath

	// CONVERT PNG to JPEG
	err = utils.ConvertToJPEG(pngFile, jpegOutputPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting PNG to JPEG: %v", err), http.StatusInternalServerError)
		return
	}

	// SERVE THE CONVERTED JPEG FILE
	http.ServeFile(w, r, jpegOutputPath)
}
