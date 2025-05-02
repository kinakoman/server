// フォルダ削除ハンドラ
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

type DeleteFolderRequest struct {
	Folder string `json:"folder"`
}

// /images/folder/delete/
type DeleteFolderHandler struct{}

func (h *DeleteFolderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var requests []DeleteFolderRequest
	// リクエストからjSONを取得

	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		http.Error(w, "Invalid Value", http.StatusBadGateway)
		return
	}

	// 画像保存先ディレクトリを取得
	originalImageFolder := os.Getenv("ORIGINAL_IMAGE_STORAGE_PATH")
	compressedImageFolder := os.Getenv("COMPRESSED_IMAGE_STORAGE_PATH")

	// データベースとの接続を確立
	con, err := connection.ConnectDB()
	if err != nil {
		http.Error(w, "DataBase is NOT running", http.StatusOK)
		return
	}
	defer con.Close()

	var deleted []*DeleteFolderRequest

	for _, req := range requests {
		// 消去対象のフォルダを取得
		targetFolder := req.Folder

		// フォルダ名のバリデーションを実行
		if module.ValidateRequestPath(originalImageFolder, targetFolder) || module.ValidateRequestPath(compressedImageFolder, targetFolder) {
			log.Println("detect invalid folder")
			continue
		}

		// 消去対象のフォルダパスを取得
		targetOriginalFolder := filepath.Join(originalImageFolder, targetFolder)
		targetCompressedFolder := filepath.Join(compressedImageFolder, targetFolder)

		// オリジナル、軽量版でフォルダを削除
		if err := os.RemoveAll(targetOriginalFolder); err != nil {
			log.Println("Remove Image Error:", targetOriginalFolder)
			continue
		}
		if err := os.RemoveAll(targetCompressedFolder); err != nil {
			log.Println("Remove Image Error:", targetCompressedFolder)
			continue
		}

		// データベースから対象のフォルダのデータ全てを削除
		if err := connection.ExecDeleteFolder(con, targetFolder); err != nil {
			log.Println("Failed to Delete Folder Info : ", targetFolder)
			continue
		}

		deleted = append(deleted, &DeleteFolderRequest{
			Folder: targetFolder,
		})

	}

	if err := json.NewEncoder(w).Encode(deleted); err != nil {
		http.Error(w, "Failed to Encode", http.StatusOK)
	}
}
