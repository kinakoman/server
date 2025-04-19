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

	for _, req := range requests {
		// 消去対象のフォルダ名とファイルを取得
		targetFolder := req.Folder
		// フォルダ名のバリデーションを実行
		if err := module.ValdateRequestPath(originalImageFolder, targetFolder); err != nil {
			http.Error(w, "invalid folder name", http.StatusBadGateway)
			return
		}
		if err := module.ValdateRequestPath(compressedImageFolder, targetFolder); err != nil {
			http.Error(w, "invalid folder name", http.StatusBadGateway)
			return
		}
		// 消去対象のファイル名を取得
		targetFilename := req.Filename
		// ファイル名のバリデーションを実行
		if err := module.ValdateRequestPath(originalImageFolder, targetFilename); err != nil {
			http.Error(w, "invalid folder name", http.StatusBadGateway)
			return
		}
		if err := module.ValdateRequestPath(compressedImageFolder, targetFilename); err != nil {
			http.Error(w, "invalid folder name", http.StatusBadGateway)
			return
		}

		// 対象ファイルのパスを取得
		targetOriginalImage := filepath.Join(originalImageFolder, targetFolder, targetFilename)
		targetCompressedImage := filepath.Join(compressedImageFolder, targetFolder, targetFilename)
		// オリジナル、軽量版の画像を削除
		if err := os.RemoveAll(targetOriginalImage); err != nil {
			log.Println("Remove Image Error:", targetOriginalImage)
		}
		if err := os.RemoveAll(targetCompressedImage); err != nil {
			log.Println("Remove Image Error:", targetCompressedImage)
		}

		// データベースから画像情報を削除
		if err := connection.ExecDelete(con, targetFolder, targetFilename); err != nil {
			log.Println("Failed to Delete Image Info : ", targetFolder, targetFilename)
		}
	}

	// ステータスコード200番を返す
	w.WriteHeader(http.StatusOK)
}
