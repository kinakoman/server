// 画像アップロードハンドラ
// 画像データをサーバ上に、画像情報をSQLに保存
// 画像データは一次保存フォルダに仮保存した後、リクエストの指定フォルダに書き換え
// 画像保存、SQLの書き込みのいずれの処理でもエラーが発生した場合は、
// 保存状態をリセット(ロールバック)し、メッセージを返す
package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"

	"image-server/connection"
	"image-server/module"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// 画像のアップロードを実行するハンドラ
type UploadHandler struct{}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// 送信された画像を保存するディレクトリ
	imageStoragePath := os.Getenv("ORIGINAL_IMAGE_STORAGE_PATH")
	// 送信された画像の軽量版の保存先ディレクトリ
	compressedStoragePath := os.Getenv("COMPRESSED_IMAGE_STORAGE_PATH")

	// 画像の一次保存先
	tempId := uuid.New().String()
	temporaryFolder := filepath.Join("temporary-", tempId)
	temporaryFolderPath := fmt.Sprintf("%s/%s", imageStoragePath, temporaryFolder)

	// 一次保存先の作成
	if err := os.MkdirAll(temporaryFolderPath, os.ModePerm); err != nil {
		http.Error(w, "Failed to Create folder to save images", http.StatusOK)
		return
	}
	defer os.RemoveAll(temporaryFolderPath) //終了時に一次保存先は削除

	// 画像データの保存先フォルダ名リスト
	var folderNameList []string
	// 保存が完了した画像ファイル
	var savedTemporary []string

	// リクエストのマルチパートのアクセスリーダーを取得
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Invalid Data", http.StatusOK)
		return
	}

	for {
		// 各パートを繰り返し実行
		part, err := reader.NextPart()
		// 最後までパートを読んだらループ終了
		if err == io.EOF {
			break
		}
		// パートの読み出しに失敗したらレスポンス
		if err != nil {
			http.Error(w, "Invalid Data", http.StatusOK)
			return
		}

		// リクエストのname属性を取得
		formName := part.FormName()

		switch formName {
		// name:folderの場合の処理⇒フォルダ名を取得
		case os.Getenv("IMAGE_SERVER_REQUEST_NAME_FOLDER"):
			// フォルダ名の読み出し
			data, err := io.ReadAll(part)
			if err != nil {
				continue
			}
			part.Close()
			// 読みだしたフォルダ名を取得
			folderNameList = append(folderNameList, string(data))

		// name:imagesの場合の処理⇒画像の保存
		case os.Getenv("IMAGE_SERVER_REQUEST_NAME_IMAGES"):
			// ファイル名が存在すれば実行
			if part.FileName() != "" {
				// ファイル名を取得
				fileName := part.FileName()
				// 一次保存フォルダに保存先を作成
				temporarySavePath := filepath.Join(temporaryFolderPath, fileName)
				temporarySavefile, err := os.Create(temporarySavePath) // 保存先ファイルの作成
				// ファイル作成エラーが発生したら処理をスキップ
				if err != nil {
					continue
				}

				// 保存先に画像データをコピー
				if _, err := io.Copy(temporarySavefile, part); err != nil {
					temporarySavefile.Close()
					part.Close()
					continue
				}

				temporarySavefile.Close()
				part.Close()

				// 保存に成功した画像ファイル名を記録
				savedTemporary = append(savedTemporary, fileName)
			} else {
				part.Close()
			}
		default:
			part.Close()
		}
	}
	// フォルダ名が複数指定されていればエラーのレスポンス
	if len(folderNameList) > 1 {
		http.Error(w, fmt.Sprintf("You cannot send multiple folders\nfolders:%s", folderNameList), http.StatusOK)
		return
	} else if len(folderNameList) == 0 || folderNameList[0] == "" { //指定が無ければフォルダ名をdefaultに
		folderNameList = append([]string{"default"}, folderNameList...)
	}

	// フォルダ名をリクエストの指定フォルダに設定
	folderName := folderNameList[0]
	// フォルダ名のバリデーションを実行
	if module.ValidateRequestPath(imageStoragePath, folderName) || module.ValidateRequestPath(compressedStoragePath, folderName) {
		http.Error(w, "invalid folder name", http.StatusBadGateway)
		return
	}

	// データベースとの接続を確立
	con, err := connection.ConnectDB()
	if err != nil {
		http.Error(w, "DataBase is NOT running", http.StatusOK)
		return
	}
	defer con.Close()

	// リクエストで指定されたフォルダの最終ディレクトリに保存された画像ファイル
	var finalSaved []string

	// 最終的な保存先ディレクトリ
	saveFolderPath := filepath.Join(imageStoragePath, folderName)

	// フォルダがすでに存在しているかチェック
	if _, err := os.Stat(saveFolderPath); os.IsNotExist(err) {
		// データベースにフォルダを追加
		if err := connection.ExecMakeFolder(con, folderName); err != nil {
			log.Println("Failed to Make Folder Info : ", folderName)
		}
		log.Println("Folder created:", saveFolderPath)
	} else if err != nil {
		// stat に失敗（権限エラーなど）
		http.Error(w, "Failed to check folder existence", http.StatusInternalServerError)
		return
	}
	// ディレクトリの作成
	if err := os.MkdirAll(saveFolderPath, os.ModePerm); err != nil {
		http.Error(w, "Failed to Create folder to save images", http.StatusOK)
		return
	}

	// 軽量版画像の保存ディレクトリ
	compressedFolderPath := filepath.Join(compressedStoragePath, folderName)
	if err := os.MkdirAll(compressedFolderPath, os.ModePerm); err != nil {
		http.Error(w, "Failed to Create folder to save images", http.StatusOK)
		return
	}

	// 一次保存が完了したファイルを繰り返し実行
	for _, fileName := range savedTemporary {
		// 一次保存先でのファイルパス
		temporarySavePath := filepath.Join(temporaryFolderPath, fileName)
		// 最終保存先でのファイルパス
		finalSavePath := filepath.Join(saveFolderPath, fileName)
		// 軽量版保存先でのファイルパス
		compressedPath := filepath.Join(compressedFolderPath, fileName)
		// 画像情報をデータベースに登録
		if err := connection.ExecSave(con, folderName, fileName); err != nil {
			continue
		}

		// ファイルを一次保存先から最終保存先に移動
		if err := os.Rename(temporarySavePath, finalSavePath); err != nil { // 保存に失敗
			// すでに同名のファイルが保存完了しているかチェック
			var flag bool
			for _, savedFilename := range finalSaved {
				if savedFilename == fileName {
					// 保存済みファイルに同名ファイルがあればフラグを立てる
					flag = true
				}
			}
			// フラグが立っていなければデータベースをロールバック
			if !flag {
				// データベースから該当ファイル情報を削除
				if err := connection.ExecDeleteLatest(con, folderName, fileName); err != nil {
					log.Println("Save failed, but the information remains in the database")
				}
			}
			continue
		}

		// 画像の軽量版を作成し保存
		if err := module.Resize(finalSavePath, compressedPath); err != nil { // 軽量化に失敗
			log.Println("Failed to compress:", fileName, err)
			// 該当ファイルを削除
			if err := os.RemoveAll(finalSavePath); err != nil {
				log.Println("Original image saving was successful, lightweight version failed:", finalSavePath)
			}
			// データベースから該当ファイル情報を削除(ロールバック)
			if err := connection.ExecDeleteLatest(con, folderName, fileName); err != nil {
				log.Println("Save failed, but the information remains in the database")
			}
			continue
		}
		// 最終保存、軽量版保存、データベース情報登録に成功した画像を記録
		finalSaved = append(finalSaved, fileName)
	}

	// レスポンス返却（JSONで保存済ファイル一覧）
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{
		"folder": {folderName},
		"saved":  finalSaved,
	})

}
