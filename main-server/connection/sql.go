package connection

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// データベースとの接続を確立する関数AC
func ConnectDB() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// データベースアクセスのDSN
	// 標準時はMySQLに合わせてJSTに
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo", dbUser, dbPassword, dbHost, dbPort, dbName)

	// データベースに接続
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err // エラーの場合nilを返す
	}

	// データベースとの接続を確認
	if err := db.Ping(); err != nil {
		return nil, err // エラーの場合nilを返す
	}

	return db, nil
}

// httpリクエストをテーブルに記録
func LogRequest(db *sql.DB, r *http.Request) {
	dbTable := os.Getenv("DB_TABLE")
	// クエリのテーブル名部分を動的に挿入
	query := fmt.Sprintf("INSERT INTO %s (path, method) VALUES (?, ?)", dbTable)
	_, err := db.Exec(query, r.URL.Path, r.Method)
	if err != nil {
		log.Println(err)
	}
}

// バックアップ記録用の軽量化ログ
func LogBackUpRequest(db *sql.DB, path string, method string, timestamp time.Time) {
	dbTable := os.Getenv("DB_TABLE")
	// クエリのテーブル名部分を動的に挿入
	query := fmt.Sprintf("INSERT INTO %s (path, method,timestamp) VALUES (?, ?, ?)", dbTable)
	_, err := db.Exec(query, path, method, timestamp)
	if err != nil {
		log.Println(err)
	}
}
