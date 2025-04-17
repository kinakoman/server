package module

import (
	"log"
	"net/http"
	"os"
)

// アクセス管理のミドルウェア(リバースプロキシサーバー以外アクセスをブロック)
func AccessCheckMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ヘッダーの秘密鍵情報を読み取る
		if r.Header.Get(os.Getenv("PROXY_SECRET_HEADER")) != os.Getenv("PROXY_SECRET_KEY") {
			// 秘密鍵が読み取れなかった場合502番エラー
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}
		// ログにメインサーバーからのアクセスを記録
		log.Println("Access from main-server")
		h.ServeHTTP(w, r)
	})
}
