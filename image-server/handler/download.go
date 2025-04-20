package handler

import (
	"encoding/json"
	"image-server/module"
	"net/http"
	"os"
	"path/filepath"
)

type DownloadRequest struct {
	Folder   string `json:"folder"`
	Filename string `json:"filename"`
}

type DownloadHandler struct{}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Value", http.StatusBadGateway)
		return
	}

	quality := r.URL.Query()["quality"]

	if len(quality) > 1 {
		http.Error(w, "You must choose one QUALITY", http.StatusOK)
	}

	imagePath := filepath.Join(req.Folder, req.Filename)

	// 画像保存先ディレクトリを取得
	originalImageFolder := os.Getenv("ORIGINAL_IMAGE_STORAGE_PATH")
	compressedImageFolder := os.Getenv("COMPRESSED_IMAGE_STORAGE_PATH")

	if module.ValdateRequestPath(originalImageFolder, imagePath) || module.ValdateRequestPath(compressedImageFolder, imagePath) {
		http.Error(w, "invalid folder name", http.StatusBadGateway)
		return
	}

	var targetImagePath string

	switch quality[0] {
	case "original":
		targetImagePath = filepath.Join(originalImageFolder, imagePath)
	case "compressed":
		targetImagePath = filepath.Join(compressedImageFolder, imagePath)
	}

	if _, err := os.Stat(targetImagePath); os.IsNotExist(err) {
		http.Error(w, "file does not exist", http.StatusOK)
		return
	}

	http.ServeFile(w, r, targetImagePath)
}
