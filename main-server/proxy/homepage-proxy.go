package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// homepageのプロキシ
func InitHomepageProxy() *httputil.ReverseProxy {
	server_scheme := os.Getenv("SERVER_SCHEME")
	homepage_host := os.Getenv("MAIN_HOMEPAGE_HOST")
	homepage_port := os.Getenv("MAIN_HOMEPAGE_PORT")

	// 接続先のurlをパース
	target, err := url.Parse(fmt.Sprintf("%s://%s:%s", server_scheme, homepage_host, homepage_port))
	if err != nil {
		log.Panicln(err)
	}
	// リバースプロキシの構造体を取得
	homepageProxy := httputil.NewSingleHostReverseProxy(target)
	// リクエストの修正を登録
	homepageProxy.Director = homepageDirector
	// エラーハンドラの登録
	homepageProxy.ErrorHandler = homepageErrorHandler

	return homepageProxy
}

// Director
func homepageDirector(r *http.Request) {
	// ヘッダーに秘密鍵情報をセット
	r.Header.Set(os.Getenv("PROXY_SECRET_HEADER"), os.Getenv("PROXY_SECRET_KEY"))
	// 接続先URLを再設定
	r.URL.Scheme = os.Getenv("SERVER_SCHEME")
	r.URL.Host = fmt.Sprintf("%s:%s", os.Getenv("MAIN_HOMEPAGE_HOST"), os.Getenv("MAIN_HOMEPAGE_PORT"))
}

// サーバー停止時のエラーハンドリング
func homepageErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("Proxy to homepage falied")
	http.Error(w, "server connection error", http.StatusBadGateway)
}
