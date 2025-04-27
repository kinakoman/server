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
			// ブラウザ以外のアクセスにはセッションが無いことをレスポンス
			if r.Method == http.MethodPost || r.URL.Path == "/csrf-token" {
				http.Error(w, "You Need New Session ID", http.StatusBadGateway)
				return
			}
			// ブラウザはログインにリダイレクト
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		//  クッキーのセッションID有効化チェック
		// セッションが無効なら/loginへリダイレクト
		if _, valid := ValidateSession(cookie.Value); !valid {
			if r.Method == http.MethodPost || r.URL.Path == "/csrf-token" {
				// ブラウザ以外のアクセスにはセッションが無いことをレスポンス
				http.Error(w, "You Need New Session ID", http.StatusBadGateway)
				return
			}
			// ブラウザはログインにリダイレクト
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// csrfトークンのチェックミドルウェア
func CsrfCheckMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストのメソッドを取得
		method := r.Method
		// POSTメソッドのみチェックを行う
		if method == http.MethodPost {
			// フォームからcsrf_tokenの値を取得

			csrfToken := r.Header.Get("csrf-token")
			// データベースで照合
			if !ValidateCsrfToken(csrfToken) {
				// 失敗したらエラー
				http.Error(w, "Invalid Post Request", http.StatusBadGateway)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
