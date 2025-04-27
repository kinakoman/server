// ログイン認証のデータベースに操作に関連した機能
// main-server/connection/connectDBの利用前提
// ユーザー情報とセッション情報の二つのデータベースを使用
package auth

import (
	"fmt"
	"log"
	"main-server/connection"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ユーザーとパスワード情報の照合　戻り値trueで照合成功、falseで失敗
func ValidateUser(username string, password string) (int, bool) {
	// データベースにアクセス
	con, err := connection.ConnectDB()
	if err != nil {
		log.Printf("\n----database error----\n%v\n----database error----\n", err)
		return 0, false
	}
	defer con.Close()

	var id int
	var hash string

	// データベースからusernameのpasswordを取得
	query := fmt.Sprintf("SELECT id ,password FROM %s WHERE username=?", os.Getenv("AUTH_USER_TABLE"))
	if err := con.QueryRow(query, username).Scan(&id, &hash); err != nil {
		// ユーザーが存在しなければ照合失敗
		return 0, false
	}

	// データベースのパスワードハッシュ値とpasswordが一致するか検証
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		// 不一致なら照合失敗
		return 0, false
	}
	// ユーザーが存在、パスワードが一致すれば照合成功
	return id, true
}

// 発行したセッションとcsrfトークンをデータベースに保存、戻り値err!=nilで保存失敗
func SaveSession(sessionID string, csrfToken string, userID int) error {
	// セッションデータベースにアクセス
	con, err := connection.ConnectDB()
	if err != nil {
		log.Printf("\n----database error----\n%v\n----database error----\n", err)
		return err
	}
	defer con.Close()

	// セッションの有効期限を発行
	expires := time.Now().Add(30 * time.Minute)

	// セッション、ユーザーID、有効期限を登録
	query := fmt.Sprintf("INSERT INTO %s (session_id,csrf_token,user_id,expires_at) values (?,?,?,?)", os.Getenv("AUTH_SESSION_TABLE"))
	_, err = con.Exec(query, sessionID, csrfToken, userID, expires)
	return err
}

// セッション情報の照合
func ValidateSession(sessionID string) (int, bool) {
	// セッションデータベースにアクセス
	con, err := connection.ConnectDB()
	if err != nil {
		log.Printf("\n----database error----\n%v\n----database error----\n", err)
		return 0, false
	}
	defer con.Close()

	var userID int
	var expires time.Time

	// セッションIDからユーザーIDと有効期限を取得
	query := fmt.Sprintf("SELECT user_id , expires_at FROM %s WHERE session_id = ?", os.Getenv("AUTH_SESSION_TABLE"))
	err = con.QueryRow(query, sessionID).Scan(&userID, &expires)
	// セッションID・ユーザーIDが存在しない、または有効期限が切れていれば照合失敗
	if err != nil || time.Now().After(expires) {
		// データベースのセッション情報を削除
		if err := DeleteSession(sessionID); err != nil {
			// セッションの削除に失敗したらサーバーに記録
			log.Printf("\n----database error----\nfailed to delete sassion\nsession_id : %s\n----database error----\n", sessionID)
		}
		// 照合失敗
		return 0, false
	}
	// 照合成功
	return userID, true
}

// CSRF情報の照合
func ValidateCsrfToken(csrfToken string) bool {
	// セッションデータベースにアクセス
	con, err := connection.ConnectDB()
	if err != nil {
		log.Printf("\n----database error----\n%v\n----database error----\n", err)
		return false
	}
	defer con.Close()

	var sessionID string

	// 送られてきたCSRFトークンが生成時のもと一致するかチェック
	query := fmt.Sprintf("SELECT session_id FROM %s WHERE csrf_token = ?", os.Getenv("AUTH_SESSION_TABLE"))
	err = con.QueryRow(query, csrfToken).Scan(&sessionID)
	// セッションID・ユーザーIDが存在しない、または有効期限が切れていれば照合失敗
	if err != nil {
		// 照合失敗
		return false
	}
	// 照合成功
	return true
}

// CSRFトークンを取得する関数
func GetCsrfToken(sessionID string) string {
	// セッションデータベースにアクセス
	con, err := connection.ConnectDB()
	if err != nil {
		log.Printf("\n----database error----\n%v\n----database error----\n", err)
		return ""
	}
	defer con.Close()
	var csrfToken string
	// セッションIDからユーザーIDと有効期限を取得
	query := fmt.Sprintf("SELECT csrf_token FROM %s WHERE session_id = ?", os.Getenv("AUTH_SESSION_TABLE"))
	err = con.QueryRow(query, sessionID).Scan(&csrfToken)
	// セッションID・ユーザーIDが存在しない、または有効期限が切れていれば照合失敗
	if err != nil {
		// 照合失敗
		return ""
	}
	return csrfToken
}

// セッションの削除、戻り値err!=nilで削除失敗
func DeleteSession(sessionID string) error {
	// セッションデータベースにアクセス
	con, err := connection.ConnectDB()
	if err != nil {
		log.Printf("\n----database error----\n%v\n----database error----\n", err)
		return err
	}
	defer con.Close()
	// セッションidを削除
	query := fmt.Sprintf("DELETE from %s WHERE session_id = ?", os.Getenv("AUTH_SESSION_TABLE"))
	_, err = con.Exec(query, sessionID)
	return err
}
