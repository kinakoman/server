package auth

import (
	"net/http"
	"os"
)

// ログイン認証ミドルウェア
func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 	リクエストのクッキーの有無のチェック
		cookie, err := r.Cookie(os.Getenv("COOKIE_SESSION_NAME"))
		// 	クッキーがなければ/loginへリダイレクト
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		//  クッキーのセッションID有効化チェック
		// セッションが無効なら/loginへリダイレクト
		if _, valid := ValidateSession(cookie.Value); !valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	})
}
