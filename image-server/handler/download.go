package handler

import (
	"image-server/module"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type DownloadHandler struct{}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	quality := r.URL.Query()["quality"]
	folder := r.URL.Query()["folder"]
	filename := r.URL.Query()["filename"]

	if len(quality) > 1 || len(folder) > 1 || len(filename) > 1 {
		http.Error(w, "You must choose one Values", http.StatusOK)
		return
	}

	imagePath := filepath.Join(folder[0], filename[0])

	log.Println(imagePath)

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
