// 画像の一覧取得に対応
// queryでフォルダを指定してフォルダ内画像情報を取得 or 全情報取得
// レスポンスはjsonの配列
package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"image-server/connection"
	"log"
	"net/http"
	"os"
	"time"
)

// レスポンスのJSON形式
type ListResponse struct {
	Id        int       `json:"id"`
	Folder    string    `json:"folder"`
	Filename  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
}

// /image/list/
type ListHandler struct{}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// リクエストのクエリパラメータからfolderの値を取得
	folderFromRequestQuery := r.URL.Query()["folder"]

	// データベースとの接続を確立
	con, err := connection.ConnectDB()
	if err != nil {
		http.Error(w, "DataBase is NOT running", http.StatusOK)
		return
	}
	defer con.Close()

	var rows *sql.Rows

	if len(folderFromRequestQuery) == 1 { // クエリパラメータが1つの場合⇒指定されたフォルダの情報を返す
		// 取得するフォルダ
		targetQuery := folderFromRequestQuery[0]
		// データベースにクエリ実行
		query := fmt.Sprintf("SELECT * FROM %s WHERE filename IS NOT NULL AND folder=?", os.Getenv("IMAGE_SERVER_NAME"))
		rows, err = con.Query(query, targetQuery)
		if err != nil {
			http.Error(w, "Failed to read Database", http.StatusOK)
			return
		}
	} else { // クエリパラメータが未指定 or 2以上⇒全情報取得
		// データベースから全データを取得
		query := fmt.Sprintf("SELECT * FROM %s WHERE filename IS NOT NULL", os.Getenv("IMAGE_SERVER_NAME"))
		rows, err = con.Query(query)
		if err != nil {
			http.Error(w, "Failed to read Database", http.StatusOK)
			return
		}
	}

	// 取得したデータを格納
	var res []*ListResponse

	// データベースから取得した各行について処理
	for rows.Next() {
		// 変数定義
		var (
			id               int
			folder, filename string
			timestamp        time.Time
		)

		// データを取得
		if err := rows.Scan(&id, &folder, &filename, &timestamp); err != nil {
			log.Println("Failed to read rows", err)
		}
		// データを配列に格納
		res = append(res, &ListResponse{
			Id:        id,
			Folder:    folder,
			Filename:  filename,
			Timestamp: timestamp,
		})
	}
	if err := rows.Close(); err != nil {
		log.Println("Failed to close rows")
	}

	// 取得してきたデータをjsonにエンコードしてレスポンス
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to Encode", http.StatusOK)
	}
}
