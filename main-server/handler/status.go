package handler

import (
	"crypto/subtle"
	"net/http"
	"os"
)

// "POST /status/"のハンドラ logに残さない代わりにLOCALからのBasicAUTHのみ受付
type StatusHandler struct{}

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 認証キーワード
	ADMIN_USER := os.Getenv("LOCAL_AUTH_USER")
	ADMIN_PASS := os.Getenv("LOCAL_AUTH_PASSWORD")
	user, pass, ok := r.BasicAuth()
	// リクエストの認証
	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(ADMIN_USER)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(ADMIN_PASS)) != 1 {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
}
