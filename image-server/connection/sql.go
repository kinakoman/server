package connection

import (
	"database/sql"
	"fmt"
	"os"

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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=Asia%%2FTokyo", dbUser, dbPassword, dbHost, dbPort, dbName)

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

// 画像情報のデータベースへの保存を実行
func ExecSave(db *sql.DB, folder string, filename string) error {
	query := fmt.Sprintf("INSERT INTO %s (folder,filename) values (?,?)", os.Getenv("IMAGE_SERVER_NAME"))
	_, err := db.Exec(query, folder, filename)
	return err
}

// 画像情報の削除を実行
func ExecDelete(db *sql.DB, folder string, filename string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE folder=? AND filename=?", os.Getenv("IMAGE_SERVER_NAME"))
	_, err := db.Exec(query, folder, filename)
	return err
}
