package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cepot-blip/opencv-go-project/utils"
)

// RESIZEREQUEST IS A STRUCTURE FOR STORING DATA IN A JSON RESIZE REQUEST
type ResizeRequest struct {
	InputPath  string `json:"input_path"`
	OutputPath string `json:"output_path"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

// RESIZE HANDLER IS A HANDLER FOR IMAGE RESIZE REQUESTSar
func ResizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// DECODE JSON REQUEST
	var req ResizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	// OPEN IMAGE FILE
	file, err := os.Open(req.InputPath)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to open image file: %v", err)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// SPICPY OUTPUT FILE
	outputPath := "assets/resized_image.jpeg"

	// RESIZE IMAGE
	err = utils.ResizeImage(file, outputPath, req.Width, req.Height)
	if err != nil {
		errorMsg := fmt.Sprintf("Error resizing image: %v", err)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	// SERVE THE RESIZED IMAGE FILE
	http.ServeFile(w, r, outputPath)
}
