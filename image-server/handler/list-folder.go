// フォルダの一覧表示に対応
package handler

import (
	"encoding/json"
	"fmt"
	"image-server/connection"
	"log"
	"net/http"
	"os"
)

// レスポンスのJSON形式
type ListFolderResponse struct {
	Folder string `json:"folder"`
}

// /image/folder/list/
type ListFolderHandler struct{}

func (h *ListFolderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// データベースとの接続を確立
	con, err := connection.ConnectDB()
	if err != nil {
		http.Error(w, "DataBase is NOT running", http.StatusOK)
		return
	}
	defer con.Close()

	// データベースから全データを取得
	query := fmt.Sprintf("SELECT DISTINCT folder FROM %s", os.Getenv("IMAGE_SERVER_NAME"))
	row, err := con.Query(query)
	if err != nil {
		http.Error(w, "Failed to read Database", http.StatusOK)
		return
	}

	// 取得したデータを格納
	var res []*ListFolderResponse

	// データベースから取得した各行について処理
	for row.Next() {
		// 変数定義
		var folder string

		// データを取得
		if err := row.Scan(&folder); err != nil {
			log.Println("Failed to read row")
		}
		// データを配列に格納
		res = append(res, &ListFolderResponse{
			Folder: folder,
		})
	}
	if err := row.Close(); err != nil {
		log.Println("Failed to close row")
	}

	// 取得してきたデータをjsonにエンコードしてレスポンス
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to Encode", http.StatusOK)
	}
}
