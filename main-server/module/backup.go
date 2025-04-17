// MySQL停止時のログのバックアップを提供するモジュール
package module

import (
	"log"
	"main-server/connection"
	"sync"
	"time"
)

// ログのバックアップ
// バックアップの構造体
type BackUpLog struct {
	BackUpUseFlag      bool                      // バックアップが呼び出されいてるかのフラグ
	BackUpRequest      []BackUpRequestContent    //リクエストのバックアップの保存先
	BackUpRequestChan  chan BackUpRequestContent // リクエストのバックアップのチャネル
	IsMySQLRunning     bool                      // MySQLの稼働ステータスフラグ
	IsReqestLogToMySQL bool                      // バックアップリクエストのデータベースへの書き込みフラグ
	mu                 sync.Mutex
}

// バックアップで保存する情報
type BackUpRequestContent struct {
	Path      string
	Method    string
	Timestamp time.Time
}

// バックアップのインスタンス生成
func NewBackUp(MaxBackUpReuestNum int) *BackUpLog {
	return &BackUpLog{
		BackUpRequest:     make([]BackUpRequestContent, 0),
		BackUpRequestChan: make(chan BackUpRequestContent, MaxBackUpReuestNum),
	}
}

// バックアップの開始関数
// 非同期で常に実行する
func (b *BackUpLog) Start() {
	// 無限ループで常に更新状態を維持
	for {
		select {
		// b.BackUpRequestからの送信を受信したら実行
		case req := <-b.BackUpRequestChan:
			// // MySQLが起動していたらバックアップの記録を実行
			// if b.IsMySQLRunning {
			// 	b.LogToMySQL()
			// 	return
			// }
			// 他の並列処理でのバックアップへのアクセスを制限
			b.mu.Lock()
			// バックアップリクエストを保存
			b.BackUpRequest = append(b.BackUpRequest, req)
			b.mu.Unlock()

		// 毎秒MySQLの稼働状態を確認
		case <-time.After(10 * time.Second):
			// MySQLが起動していたらバックアップの記録を実行
			if b.IsMySQLRunning {
				b.LogToMySQL()
			}
		}
	}
}

func (b *BackUpLog) LogToMySQL() {
	// 他の並列処理でのバックアップへのアクセスを制限
	b.mu.Lock()
	// すでにMySQLへの書き込みが進行中、またはバックアップリクエストが無ければ直ちに終了
	if b.IsReqestLogToMySQL || len(b.BackUpRequest) == 0 {
		b.mu.Unlock()
		return
	}
	// バックアップの書き込み処理を開始
	b.IsReqestLogToMySQL = true
	// バックアップリクエストのコピーを取得
	SaveBackupRequest := make([]BackUpRequestContent, len(b.BackUpRequest))
	copy(SaveBackupRequest, b.BackUpRequest)
	// バックアップリクエストを初期化
	b.BackUpRequest = b.BackUpRequest[:0]
	// バックアップへのアクセスをを解放
	b.mu.Unlock()

	// データベースにログイン
	db, err := connection.ConnectDB()
	if err != nil { // ログインエラー発生時
		log.Printf("\n----database error----\n%v\n----database error----\n", err)
	} else {
		// ログインエラーが発生しなければデータベースにログを記録
		defer db.Close()
		// 保存したバックアップリクエストをデータベースに書き込み
		for _, req := range SaveBackupRequest {
			connection.LogBackUpRequest(db, req.Path, req.Method, req.Timestamp)
		}
	}
	// バックアップの書き込み完了
	b.mu.Lock()
	b.IsReqestLogToMySQL = false
	b.mu.Unlock()
}

// バックアップリクエストを取得して送信
func (b *BackUpLog) SendBackUpReuest(req BackUpRequestContent) {
	b.BackUpRequestChan <- req
}

// MySQLの停止をバックアップに通知
func (b *BackUpLog) MySQLIsStop() {
	b.mu.Lock()
	b.IsMySQLRunning = false
	b.mu.Unlock()
}

// MySQLの起動をバックアップに通知
func (b *BackUpLog) MySQLIsRunning() {
	b.mu.Lock()
	b.IsMySQLRunning = true
	b.mu.Unlock()
}
