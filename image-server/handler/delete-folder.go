// フォルダ削除ハンドラ
package handler

import (
	"encoding/json"
	"image-server/connection"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type DeleteFolderRequest struct {
	Folder string `json:"folder"`
}

// /images/folder/delete/
type DeleteFolderHandler struct{}

func (h *DeleteFolderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// リクエストからjSONを取得
	var req DeleteFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Value", http.StatusBadGateway)
		return
	}
	// 消去対象のフォルダを取得
	targetFolder := req.Folder
	// 画像保存先ディレクトリを取得
	originalImageFolder := os.Getenv("ORIGINAL_IMAGE_STORAGE_PATH")
	compressedImageFolder := os.Getenv("COMPRESSED_IMAGE_STRAGE_PATH")
	// 消去対象のフォルダパスを取得
	targetOriginalFolder := filepath.Join(originalImageFolder, targetFolder)
	targetCompressedFolder := filepath.Join(compressedImageFolder, targetFolder)

	// オリジナル、軽量版でフォルダを削除
	if err := os.RemoveAll(targetOriginalFolder); err != nil {
		log.Println("Remove Image Error:", targetOriginalFolder)
	}
	if err := os.RemoveAll(targetCompressedFolder); err != nil {
		log.Println("Remove Image Error:", targetCompressedFolder)
	}

	// データベースとの接続を確立
	con, err := connection.ConnectDB()
	if err != nil {
		http.Error(w, "DataBase is NOT running", http.StatusOK)
		return
	}
	defer con.Close()
	// データベースから対象のフォルダのデータ全てを削除
	if err := connection.ExecDeleteFolder(con, targetFolder); err != nil {
		log.Println("Failed to Delete Folder Info : ", targetFolder)
	}
}
