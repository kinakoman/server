package handler

import (
	"encoding/json"
	"image-server/connection"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type MakeFolderRequest struct {
	Folder string `json:"folder"`
}

type MakeFolderHandler struct{}

func (h *MakeFolderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req MakeFolderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invaled Value", http.StatusBadGateway)
		return
	}

	targetFolder := req.Folder
	// 画像保存先ディレクトリを取得
	originalImageFolder := os.Getenv("ORIGINAL_IMAGE_STORAGE_PATH")
	compressedImageFolder := os.Getenv("COMPRESSED_IMAGE_STRAGE_PATH")
	// 消去対象のフォルダパスを取得
	targetOriginalFolder := filepath.Join(originalImageFolder, targetFolder)
	targetCompressedFolder := filepath.Join(compressedImageFolder, targetFolder)
	if err := os.MkdirAll(targetOriginalFolder, os.ModePerm); err != nil {
		log.Println("Failed to Create Folder", targetOriginalFolder)
	}
	if err := os.MkdirAll(targetCompressedFolder, os.ModePerm); err != nil {
		log.Println("Failed to Create Folder", targetCompressedFolder)
	}

	// データベースとの接続を確立
	con, err := connection.ConnectDB()
	if err != nil {
		http.Error(w, "DataBase is NOT running", http.StatusOK)
		return
	}
	defer con.Close()
	// データベースから対象のフォルダのデータ全てを削除
	if err := connection.ExecMakeFolder(con, targetFolder); err != nil {
		log.Println("Failed to Make Folder Info : ", targetFolder)
	}
}
