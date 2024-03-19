# TUTORIAL FOR COMPRESSED IMAGE PROCESSING IN GOLANG USING OPENCV4

### **Requirements for use**

• Go v1.22.1+

## STEPS TO MAKE COMPRESSED IMAGE PROCESSING IN GOLANG

1. Create a project folder named opencv-go-project with syntax :

   ```bash
   mkdir opencv-go-project
   ```

2. Initialize the Project by doing the following syntax :

   ```bash
   go mod init github.com/{username git}/{name folder}

   example :

   go mod init github.com/cepot-blip/opencv-go-project
   ```

3. Installing several dependencies for our project needs, here I use Gorilla mux for the Router, the syntax is as follows :

   ```bash
   go get github.com/gorilla/mux
   ```

4. install opencv to process images

   ```bash
   go get gocv.io/x/gocv
   ```

5. We also have to install opencv via home brew to use its features.

   ```bash
   brew install opencv

   don't forget to set the PATH on each of your computers
   ```

   **After completing all the dependencies we need, we go straight into our favorite code editor Visual Studio Code. and we will create 3 main folders and one file with the names assets, handlers, utils and main.go**

   ![Screenshot 2024-03-19 at 13.49.02.png](TUTORIAL%20FOR%20COMPRESSED%20IMAGE%20PROCESSING%20IN%20GOLANG%20f06a8ed97479477ca77e6ace8922841b/Screenshot_2024-03-19_at_13.49.02.png)

6. next we will start coding in the utils folder by creating 3 files named

- compressImage.go
- convertToJPEG.go
- resizeImage.go

**compressImage.go**

```go
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
```

**convertToJPEG.go**

```go
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
```

**resizeImage.go**

```go
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
```

1. Next we will create 3 files inside the handler folder with the name

- compressHandler.go
- convertHandler.go
- resizeHandler.go

**compressHandler.go**

```go
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
```

**convertHandler.go**

```go
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
```

**resizeHandler.go**

```go
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
```

1. ext we will do the coding in fila main.go to do our router HTTP server

**main.go**

```go
package main

import (
	"log"
	"net/http"

	handlers "github.com/cepot-blip/opencv-go-project/handler"
	"github.com/gorilla/mux"
)

func main() {
	// INITIALIZE THE ROUTER USING MUX.NEWROUTER()
	router := mux.NewRouter()
	router.HandleFunc("/convert", handlers.ConvertHandler).Methods("POST")
	router.HandleFunc("/resize", handlers.ResizeHandler).Methods("POST")
	router.HandleFunc("/compress", handlers.CompressHandler).Methods("POST")

	// SPECIFY THE SERVER ADDRESS
	serverAddr := ":8080"

	// PRINT A MESSAGE TO THE CONSOLE THAT THE SERVER IS STARTING AND LISTENING ON A SPECIFIC ADDRESS
	log.Printf("Server is starting and listening on %s...", serverAddr)

	// USE HTTP.LISTENANDSERVE() TO START THE SERVER WITH THE PREVIOUSLY DEFINED ROUTER
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
```

1. **next we put 1 image file into the assets folder with type png**

Then this is how it looks now

![Screenshot 2024-03-19 at 14.25.00.png](TUTORIAL%20FOR%20COMPRESSED%20IMAGE%20PROCESSING%20IN%20GOLANG%20f06a8ed97479477ca77e6ace8922841b/Screenshot_2024-03-19_at_14.25.00.png)

1. the next step we will run the server in the terminal with the following sytax :

```go
go run main.go
```

and when successful it will show a display like this on your computer terminal

![Screenshot 2024-03-19 at 14.28.41.png](TUTORIAL%20FOR%20COMPRESSED%20IMAGE%20PROCESSING%20IN%20GOLANG%20f06a8ed97479477ca77e6ace8922841b/Screenshot_2024-03-19_at_14.28.41.png)

1. okay next step we will do a test through postman

the first step is to create a new collection in postman and add a new request with the post method and this is an example of the JSON body

![Screenshot 2024-03-19 at 14.36.04.png](TUTORIAL%20FOR%20COMPRESSED%20IMAGE%20PROCESSING%20IN%20GOLANG%20f06a8ed97479477ca77e6ace8922841b/Screenshot_2024-03-19_at_14.36.04.png)

and this is the result of the successful body response

![Screenshot 2024-03-19 at 14.38.03.png](TUTORIAL%20FOR%20COMPRESSED%20IMAGE%20PROCESSING%20IN%20GOLANG%20f06a8ed97479477ca77e6ace8922841b/Screenshot_2024-03-19_at_14.38.03.png)

This is the result of the request body from convert

![Screenshot 2024-03-19 at 14.39.22.png](TUTORIAL%20FOR%20COMPRESSED%20IMAGE%20PROCESSING%20IN%20GOLANG%20f06a8ed97479477ca77e6ace8922841b/Screenshot_2024-03-19_at_14.39.22.png)

This is the result of the request body from resize

![Screenshot 2024-03-19 at 14.39.39.png](TUTORIAL%20FOR%20COMPRESSED%20IMAGE%20PROCESSING%20IN%20GOLANG%20f06a8ed97479477ca77e6ace8922841b/Screenshot_2024-03-19_at_14.39.39.png)

when finished testing in postman, the results will be automatically saved in the assets folder like this

![Screenshot 2024-03-19 at 15.57.40.png](TUTORIAL%20FOR%20COMPRESSED%20IMAGE%20PROCESSING%20IN%20GOLANG%20f06a8ed97479477ca77e6ace8922841b/Screenshot_2024-03-19_at_15.57.40.png)

And we did it
Thanks!
