// 画像のフォルダ移動ハンドラ
package handler

import (
	"encoding/json"
	"image-server/connection"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type MoveFolderRequest struct {
	PreFolder  string `json:"prefolder"`
	PostFolder string `json:"postfolder"`
	File       []struct {
		Id       int
		Filename string
	} `json:"file"`
}

type MoveFolderHandler struct{}

func (h *MoveFolderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// リクエストの取得
	var requests MoveFolderRequest
	// リクエストをjson形式で取得
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		http.Error(w, "Invalid Value", http.StatusBadGateway)
		return
	}

	// 画像保存先ディレクトリを取得
	originalImageFolder := os.Getenv("ORIGINAL_IMAGE_STORAGE_PATH")
	compressedImageFolder := os.Getenv("COMPRESSED_IMAGE_STORAGE_PATH")
	// 移動前のフォルダ名
	prefolder := requests.PreFolder
	// 移動前のフォルダパス
	preOriginalFolder := filepath.Join(originalImageFolder, prefolder)
	preCompressedFolder := filepath.Join(compressedImageFolder, prefolder)

	// 移動先のフォルダ名
	postFolder := requests.PostFolder
	// 移動先のフォルダパス
	postOriginalFolder := filepath.Join(originalImageFolder, postFolder)
	postCompressedFolder := filepath.Join(compressedImageFolder, postFolder)

	// データベースとの接続を確立
	con, err := connection.ConnectDB()
	if err != nil {
		http.Error(w, "DataBase is NOT running", http.StatusOK)
		return
	}
	defer con.Close()

	// 移動させる画像の配列
	targetFiles := requests.File
	// 移動完了したファイル
	var moved []string

	// 各ファイルに対して処理を実行
	for _, file := range targetFiles {
		// 対象画像のId
		targetId := file.Id
		// 対象画像の名前
		targetFilename := file.Filename

		// データベースの書き換え処理
		if err := connection.ExecMoveFolder(con, postFolder, targetId); err != nil {
			log.Println("Failed to Move Image Info:", targetFilename)
			continue
		}

		// 移動前の画像データパス
		preOriginalPath := filepath.Join(preOriginalFolder, targetFilename)
		preCompressedPath := filepath.Join(preCompressedFolder, targetFilename)
		// 移動後の画像データパス
		postOriginalPath := filepath.Join(postOriginalFolder, targetFilename)
		postCompressedPath := filepath.Join(postCompressedFolder, targetFilename)

		// 画像ファイルの移動処理
		// ロールバック処理未実装
		// オリジナル画像
		if err := os.Rename(preOriginalPath, postOriginalPath); err != nil {
			log.Printf("Failed to move original image:%s\n%s to %s\n", targetFilename, prefolder, postFolder)
			continue
		}
		// 軽量版画像
		if err := os.Rename(preCompressedPath, postCompressedPath); err != nil {
			log.Printf("Failed to move original image:%s\n%s to %s\n", targetFilename, prefolder, postFolder)
			continue
		}

		// 移動完了した画像を記録
		moved = append(moved, targetFilename)

	}

}
