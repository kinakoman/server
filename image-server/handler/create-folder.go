package handler

import (
	"encoding/json"
	"image-server/connection"
	"image-server/module"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type MakeFolderRequest struct {
	Folder string `json:"folder"`
}

type CreateFolderHandler struct{}

func (h *CreateFolderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req MakeFolderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, "サーバー側のエラーが発生しました。", http.StatusBadGateway)
		return
	}
	// 作成するフォルダ名を取得
	targetFolder := req.Folder

	// 画像保存先ディレクトリを取得
	originalImageFolder := os.Getenv("ORIGINAL_IMAGE_STORAGE_PATH")
	compressedImageFolder := os.Getenv("COMPRESSED_IMAGE_STORAGE_PATH")

	// フォルダ名のバリデーションを実行
	if module.ValdateRequestPath(originalImageFolder, targetFolder) || module.ValdateRequestPath(compressedImageFolder, targetFolder) {
		http.Error(w, "フォルダ名に不正な文字が使用されています。", http.StatusBadGateway)
		return
	}

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
		http.Error(w, "データベースが停止しています。", http.StatusOK)
		return
	}
	defer con.Close()
	// データベースにフォルダ情報を登録
	if err := connection.ExecMakeFolder(con, targetFolder); err != nil {
		log.Println("Failed to Make Folder Info : ", targetFolder)
	}
}
