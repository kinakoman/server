// 画像ダウンロード
// /image/download?folder=フォルダ名&filename=ファイル名&quality=品質設定でアクセス
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
	// クエリーからパラメータを取得
	folder := r.URL.Query()["folder"]
	filename := r.URL.Query()["filename"]
	quality := r.URL.Query()["quality"]

	// クエリーが1つ以上のリクエストがあればエラー
	if len(quality) > 1 || len(folder) > 1 || len(filename) > 1 {
		http.Error(w, "You must choose one Values", http.StatusOK)
		return
	}

	// 画像のパスを作成
	imagePath := filepath.Join(folder[0], filename[0])

	// 画像保存先ディレクトリを取得
	originalImageFolder := os.Getenv("ORIGINAL_IMAGE_STORAGE_PATH")
	compressedImageFolder := os.Getenv("COMPRESSED_IMAGE_STORAGE_PATH")

	// 画像パスのバリデーション
	if module.ValidateRequestPath(originalImageFolder, imagePath) || module.ValidateRequestPath(compressedImageFolder, imagePath) {
		log.Println("now")
		http.Error(w, "invalid folder name", http.StatusBadGateway)
		return
	}

	// 最終的な取得先の画像パス
	var targetImagePath string

	// qualityクエリーの値に応じて参照先画像を変更
	switch quality[0] {
	case "original": // origanl画像
		targetImagePath = filepath.Join(originalImageFolder, imagePath)
	case "compressed": // 軽量画像
		targetImagePath = filepath.Join(compressedImageFolder, imagePath)
	}

	// ファイルの存在確認
	if _, err := os.Stat(targetImagePath); os.IsNotExist(err) {
		http.Error(w, "file does not exist", http.StatusOK)
		return
	}

	// 画像ファイルをレスポンス
	http.ServeFile(w, r, targetImagePath)
}
