package module

import (
	"log"
	"main-server/auth"
	"main-server/connection"
	"net/http"
	"time"
)

// ミドルウェアの一括適用
func SetMiddleware(h http.Handler, Backup *BackUpLog) http.Handler {
	return LogMiddleware(auth.AuthMiddleware(h), Backup)
}

// ログミドルウェア
// リクエストをデータベースに書き込む
// リクエストのバックアップを保存する
func LogMiddleware(h http.Handler, BackUp *BackUpLog) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// サーバー側にアクセスログを出力
		log.Printf("%v => %v", r.Method, r.URL.Path)

		// データベースにログイン
		db, err := connection.ConnectDB()
		if err != nil { // ログインエラー発生時
			log.Printf("\n----database error----\n%v\n----database error----\n", err)
			// バックアップにデータベースの停止を通知
			BackUp.MySQLIsStop()
			BackUp.SendBackUpReuest(BackUpRequestContent{Path: r.URL.Path, Method: r.Method, Timestamp: time.Now()})
		} else {
			// ログインエラーが発生しなければデータベースにログを記録
			defer db.Close()
			connection.LogRequest(db, r)
			// バックアップにデータベースの起動を通知
			if !BackUp.IsMySQLRunning {
				BackUp.MySQLIsRunning()
			}
		}
		h.ServeHTTP(w, r)
	})
}
