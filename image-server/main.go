package main

import (
	"fmt"
	"image-server/handler"
	"image-server/module"
	"log"
	"net/http"

	"os"

	"github.com/joho/godotenv"
)

// ルーティングの初期化関数
func initRoute() *http.ServeMux {
	mux := http.NewServeMux()

	// ルーティングの設定
	mux.Handle("/", module.AccessCheckMiddleware(&handler.RootHandler{}))
	// 画像のアップロード
	mux.Handle("POST /image/upload/", module.AccessCheckMiddleware(&handler.UploadHandler{}))
	// 画像リスト取得
	mux.Handle("GET /image/list", module.AccessCheckMiddleware(&handler.ListHandler{}))
	// 画像削除
	mux.Handle("POST /image/delete/", module.AccessCheckMiddleware(&handler.DeleteHandler{}))
	// フォルダの削除
	mux.Handle("POST /image/folder/delete/", module.AccessCheckMiddleware(&handler.DeleteFolderHandler{}))
	// フォルダの作成
	mux.Handle("POST /image/folder/create/", module.AccessCheckMiddleware(&handler.CreateFolderHandler{}))
	// フォルダリスト取得
	mux.Handle("GET /image/folder/list/", module.AccessCheckMiddleware(&handler.ListFolderHandler{}))
	return mux

}

func main() {
	// .envファイルのロード
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	HOST := os.Getenv("IMAGE_SERVER_HOST")
	PORT := os.Getenv("IMAGE_SERVER_PORT")
	// サーバーアドレス
	Addr := fmt.Sprintf("%s:%s", HOST, PORT)

	// マルチプレクサの初期化
	mux := initRoute()
	// サーバーの構造体の初期化
	server := http.Server{Addr: Addr, Handler: mux}

	// サーバーの起動
	log.Printf("Activate Server")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
