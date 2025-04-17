// 画像の一覧取得に対応
package handler

import (
	"encoding/json"
	"fmt"
	"image-server/connection"
	"log"
	"net/http"
	"os"
	"time"
)

type ListResponse struct {
	Id        int       `json:"id"`
	Folder    string    `json:"folder"`
	Filename  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
}

type ListHandler struct{}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// データベースとの接続を確立
	con, err := connection.ConnectDB()
	if err != nil {
		http.Error(w, "DataBase is NOT running", http.StatusOK)
		return
	}
	defer con.Close()

	query := fmt.Sprintf("SELECT * FROM %s", os.Getenv("IMAGE_SERVER_NAME"))
	row, err := con.Query(query)
	if err != nil {
		http.Error(w, "Failed to read Database", http.StatusOK)
		return
	}

	var res []*ListResponse

	for row.Next() {
		var (
			id               int
			folder, filename string
			timestamp        time.Time
		)

		if err := row.Scan(&id, &folder, &filename, &timestamp); err != nil {
			log.Println("Failed to read row")
		}
		log.Println(id)

		res = append(res, &ListResponse{
			Id:        id,
			Folder:    folder,
			Filename:  filename,
			Timestamp: timestamp,
		})
	}
	if err := row.Close(); err != nil {
		log.Println("Failed to close row")
	}
	log.Println(res)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to Encode", http.StatusOK)
	}
}
