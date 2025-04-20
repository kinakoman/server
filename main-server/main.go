package main

import (
	"fmt"
	"log"
	"main-server/auth"
	"main-server/handler"
	"main-server/module"
	"main-server/proxy"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// ルーティングの初期化関数
func initRoute(BackUp *module.BackUpLog) *http.ServeMux {
	// マルチプレクサの開始
	mux := http.NewServeMux()

	// ルートハンドラ
	homepageProxy := proxy.InitHomepageProxy() // homepage-serverのプロキシ
	mux.Handle("/", module.SetMiddleware(homepageProxy, BackUp))

	// test用ルート
	mux.Handle("/test/", auth.AuthMiddleware(&handler.RootHandler{}))

	// imageハンドラ
	imageProxy := proxy.InitImageProxy() // image-serverのプロキシ
	mux.Handle("/image/", module.SetMiddleware(imageProxy, BackUp))

	// ログインハンドラ
	mux.Handle("/login", &auth.LoginHandler{})
	mux.Handle("/logout", &auth.LogoutHandler{})

	// サーバーモニターへのステータス変更ハンドラ
	mux.Handle("POST /status/", &handler.StatusHandler{})
	return mux
}

func main() {
	// .envファイルのロード
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	HOST := os.Getenv("MAIN_SERVER_HOST")
	PORT := os.Getenv("MAIN_SERVER_PORT")
	// サーバーアドレス
	Addr := fmt.Sprintf("%s:%s", HOST, PORT)

	// バックアップのインスタンス化
	BackUp := module.NewBackUp(50)
	// バックアップ管理を非同期で開始
	go BackUp.Start()

	// マルチプレクサを初期化
	mux := initRoute(BackUp)
	// サーバー構造体の初期化
	server := http.Server{Addr: Addr, Handler: mux}
	// サーバーの起動
	log.Printf("Activate Server")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
