// 画像の消去ハンドラ
// リクエストのjsonは配列形式
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

// deleteリクエストのJSON
type DeleteRequestStruct struct {
	Id       int    `json:"id"`
	Folder   string `json:"folder"`
	Filename string `json:"filename"`
}

// /delete/
type DeleteHandler struct{}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// リクエストからjsonを取得
	var requests []DeleteRequestStruct
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

	var deleted []*DeleteRequestStruct

	for _, req := range requests {
		// 消去対象のフォルダ名とファイルを取得
		targetFolder := req.Folder
		// フォルダ名のバリデーションを実行
		if module.ValdateRequestPath(originalImageFolder, targetFolder) || module.ValdateRequestPath(compressedImageFolder, targetFolder) {
			log.Println("detect invalid file")
			continue
		}

		// 消去対象のファイル名を取得
		targetFilename := req.Filename
		// ファイル名のバリデーションを実行
		if module.ValdateRequestPath(originalImageFolder, targetFilename) || module.ValdateRequestPath(compressedImageFolder, targetFilename) {
			log.Println("detect invalid file")
			continue
		}

		// 消去対象のidを取得
		targetId := req.Id

		// データベースから画像情報を削除
		if err := connection.ExecDelete(con, targetId, targetFolder, targetFilename); err != nil {
			log.Println("Failed to Delete Image Info : ", targetFolder, targetFilename)
			continue
		}

		if connection.ImageDataNoExist(con, targetFolder, targetFilename) {
			// 対象ファイルのパスを取得
			targetOriginalImage := filepath.Join(originalImageFolder, targetFolder, targetFilename)
			targetCompressedImage := filepath.Join(compressedImageFolder, targetFolder, targetFilename)
			// オリジナル、軽量版の画像を削除
			if err := os.RemoveAll(targetOriginalImage); err != nil {
				log.Println("Remove Image Error:", targetOriginalImage)
				continue
			}
			if err := os.RemoveAll(targetCompressedImage); err != nil {
				log.Println("Remove Image Error:", targetCompressedImage)
				continue
			}
		}

		deleted = append(deleted, &DeleteRequestStruct{
			Folder:   targetFolder,
			Filename: targetFilename,
		})
	}
	// 取得してきたデータをjsonにエンコードしてレスポンス
	if err := json.NewEncoder(w).Encode(deleted); err != nil {
		http.Error(w, "Failed to Encode", http.StatusOK)
	}
}
