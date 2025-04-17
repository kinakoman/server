package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// imageプロキシの初期化関数
func InitImageProxy() *httputil.ReverseProxy {
	server_scheme := os.Getenv("SERVER_SCHEME")
	image_server_host := os.Getenv("IMAGE_SERVER_HOST")
	image_server_port := os.Getenv("IMAGE_SERVER_PORT")

	// 接続先のurlをパース
	target, err := url.Parse(fmt.Sprintf("%s://%s:%s", server_scheme, image_server_host, image_server_port))
	if err != nil {
		log.Println(err)
	}
	// リバースプロキシの構造体を取得
	imageProxy := httputil.NewSingleHostReverseProxy(target)
	// リクエストの修正を登録
	imageProxy.Director = imageDirector
	// レスポンス修正を登録
	imageProxy.ModifyResponse = imageModifyResponse
	// エラーハンドラを登録
	imageProxy.ErrorHandler = imageErrorHandler

	return imageProxy
}

// Director
func imageDirector(r *http.Request) {
	// ヘッダーに秘密鍵を登録
	r.Header.Set(os.Getenv("PROXY_SECRET_HEADER"), os.Getenv("PROXY_SECRET_KEY"))
	// 接続先URLを再設定
	r.URL.Scheme = os.Getenv("SERVER_SCHEME")
	r.URL.Host = fmt.Sprintf("%s:%s", os.Getenv("IMAGE_SERVER_HOST"), os.Getenv("IMAGE_SERVER_PORT"))
}

// ModifyResponse
func imageModifyResponse(r *http.Response) error {
	// ステータス200番以外は処理を実行
	if r.StatusCode != http.StatusOK {
		defer r.Body.Close()

		// ステータスを502番に書き換え
		r.StatusCode = http.StatusBadGateway
		newBody := fmt.Sprintf("%s\n", http.StatusText(http.StatusBadGateway))
		r.Body = io.NopCloser(strings.NewReader(newBody))
		r.Header.Set("Content-Length", strconv.Itoa(len(newBody)))
		r.TransferEncoding = nil
	}
	return nil
}

// 停止時のエラーハンドリング
func imageErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("Proxy to image-server falied")
	http.Error(w, "server connection error", http.StatusBadGateway)
}
