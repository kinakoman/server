// 入力されたパスのバリデーション
// パスの値そのものに加えて、basePath(配下にパスを置きたいパス)以下に配置されているかチェック
// エラー!=nilで入力されたパスは不正
package module

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ValdateRequestPath(basePath string, request string) error {
	// パスが空白の場合
	if request == "" {
		return errors.New("value is empty")
	}
	// パスの長さ100以上の時
	if len(request) > 100 {
		return errors.New("value is too long")
	}
	// パスに不正な文字が使われている
	invalidChars := regexp.MustCompile(`[<>:"\\|?*\x00-\x1F]`)
	if invalidChars.MatchString(request) {
		return errors.New("value contain invalid character")
	}
	// basePathとrequestの結合
	fullPath := filepath.Join(basePath, request)
	// 相対パスなどを正規化
	cleanPath := filepath.Clean(fullPath)
	// basePathの正規化
	basePathClean := filepath.Clean(basePath) + string(os.PathSeparator)
	// requestがbasePathの配下になっているかチェック
	if !strings.HasPrefix(cleanPath, basePathClean) {
		return errors.New("value contains escape")
	}

	return nil
}
