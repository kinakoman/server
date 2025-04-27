// ログイン認証関係のハンドラ
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

// ログイン画面のテンプレート
var loginTmpl = template.Must(template.ParseFiles("template/login.html"))

// /login ハンドラ
type LoginHandler struct{}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ゲットメソッドならログイン画面の表示
	if r.Method == http.MethodGet {
		loginTmpl.Execute(w, nil)
		return
	}
	// リクエストからusernameとpassword取得
	username := r.FormValue("username")
	password := r.FormValue("password")

	// usernameとpasswordがデータベースに保存されいてるかチェック
	userID, valid := ValidateUser(username, password)
	if !valid { // 保存されていなければログイン画面に戻す
		loginTmpl.Execute(w, "Invalid credentials")
		return
	}

	// セッションIDの生成
	sessionID, csrfToken, err := GenerateSessionID()
	if err != nil {
		log.Println("Failed to start session:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// セッションIDの登録
	if SaveSession(sessionID, csrfToken, userID) != nil {
		log.Println("Failed to save session:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// リクエストに付与するクッキーの設定
	cookie := &http.Cookie{
		Name:     os.Getenv("COOKIE_SESSION_NAME"),
		Value:    sessionID,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode, // クロスサイト送信拒否
		Path:     "/",
		// Secure: true,
		Expires: time.Now().Add(10 * time.Second), // 有効期限は30分
	}
	// クッキーのセット
	http.SetCookie(w, cookie)
	// /にリダイレクト
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// セッションIDとcsrfトークンの生成関数
func GenerateSessionID() (string, string, error) {
	sessionID := make([]byte, 24)
	csrfToken := make([]byte, 32)
	_, err := rand.Read(sessionID)
	if err != nil {
		return "", "", err
	}
	_, err = rand.Read(csrfToken)
	if err != nil {
		return "", "", err
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sessionID), base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(csrfToken), nil
}

// /crsf-token CSRFトークンを取得するハンドラ
type GetCsrfTokenHandler struct{}

func (h GetCsrfTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 	リクエストのクッキーの有無のチェック
	cookie, err := r.Cookie(os.Getenv("COOKIE_SESSION_NAME"))
	// 	クッキーがなければ/loginへリダイレクト
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//  クッキーのセッションID有効化チェック
	sessionID := cookie.Value

	csrfToken := GetCsrfToken(sessionID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": csrfToken,
	})
}

// /logout ハンドラ
type LogoutHandler struct{}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// リクエストのクッキーからセッションIDを取得
	cookie, err := r.Cookie(os.Getenv("COOKIE_SESSION_NAME"))
	if err == nil {
		// セッションIDの削除
		if err := DeleteSession(cookie.Value); err != nil {
			// セッションの削除に失敗したらサーバーに記録
			log.Printf("\n----database error----\nfailed to delete sassion\nsession_id : %s\n----database error----\n", cookie.Value)
		}
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
