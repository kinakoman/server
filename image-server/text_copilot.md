# 写真保存APIサーバー

このサーバーは、HTTP リクエストを通じて画像ファイルのアップロード、取得、削除、移動、およびフォルダ管理を行います。画像自体はサーバー内のimagesディレクトリに保存され、画像情報はSQLデータベースに登録されます。

## 目次

- [概要](#概要)
- [セットアップと起動方法](#セットアップと起動方法)
- [環境変数](#環境変数)
- [API エンドポイント](#api-エンドポイント)
  - [POST /image/upload/](#post-imageupload)
  - [GET /image/list](#get-imagelist)
  - [POST /image/delete/](#post-imagedelete)
  - [POST /image/move/](#post-imagemove)
  - [POST /image/folder/create/](#post-imagefoldercreate)
  - [POST /image/folder/delete/](#post-imagefolderdelete)
  - [GET /image/download](#get-imagedownload)
- [依存関係](#依存関係)

## 概要

- 画像ファイルは **original** と **compressed** の２種類のディレクトリに保存されます。  
- アップロード時に一時保存ディレクトリを使い、アップロード完了後に指定フォルダへ移動および軽量版画像を生成します。  
- アクセスは[リバースプロキシサーバー](module/middleware.go)からのみ許可されています（ヘッダーによる認証）。

## セットアップと起動方法

1. 必要なGoのバージョンと依存ライブラリがインストールされていることを確認してください。

2. 環境変数を設定するために、プロジェクトルートの上位ディレクトリに `.env` ファイルを用意します。  
   例:
   ```env
   IMAGE_SERVER_HOST=localhost
   IMAGE_SERVER_PORT=8080
   ORIGINAL_IMAGE_STORAGE_PATH=./images/original
   COMPRESSED_IMAGE_STORAGE_PATH=./images/compressed
   IMAGE_SERVER_REQUEST_NAME_FOLDER=folder
   IMAGE_SERVER_REQUEST_NAME_IMAGES=images
   IMAGE_SERVER_NAME=your_table_name

   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=your_db_name

   PROXY_SECRET_HEADER=X-Proxy-Secret
   PROXY_SECRET_KEY=your_secret_key
   ```

3. アプリケーションをビルドし、起動します:
   ```sh
   go mod tidy
   go run main.go
   ```

## 環境変数

- **IMAGE_SERVER_HOST**: サーバーのホスト名  
- **IMAGE_SERVER_PORT**: サーバーのポート番号  
- **ORIGINAL_IMAGE_STORAGE_PATH**: オリジナル画像の保存先ディレクトリパス  
- **COMPRESSED_IMAGE_STORAGE_PATH**: 軽量版画像の保存先ディレクトリパス  
- **IMAGE_SERVER_REQUEST_NAME_FOLDER**: アップロード時のform-dataにおけるフォルダ名パラメータのname  
- **IMAGE_SERVER_REQUEST_NAME_IMAGES**: アップロード時のform-dataにおける画像ファイルのname  
- **IMAGE_SERVER_NAME**: SQLデータベース上のテーブル名  

- **DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME**: データベース接続情報  
- **PROXY_SECRET_HEADER, PROXY_SECRET_KEY**: リバースプロキシ認証のためのヘッダー情報  

## API エンドポイント

### POST /image/upload/
**説明:**  
画像ファイルのアップロードおよびデータベースへの登録、軽量版画像の生成を行います。  
リクエストの内容により、一時的に画像を保存し、指定のフォルダへ移動します。  
アップロード処理は[handler/upload.go](handler/upload.go)で実装されています。

**リクエスト形式:**  
- **Content-Type:** multipart/form-data  
- **パラメータ:**
  - **folder** (テキスト)  
    保存先のフォルダ名（指定がない場合は `"default"` として扱われます）。
  - **images** (ファイル)  
    複数可。アップロード対象の画像ファイル。

**レスポンス (JSON):**
```json
{
  "folder": ["保存先フォルダ名"],
  "saved": ["保存が完了した画像ファイル名のリスト"]
}
```

---

### GET /image/list
**説明:**  
アップロードされた画像情報を一覧取得します。  
クエリパラメータでフォルダを指定すると、そのフォルダ内の画像情報のみ取得します。  
詳細は[handler/list.go](handler/list.go)を参照ください。

**リクエストパラメータ:**
- **folder** (optional, クエリパラメータ)  
  取得対象フォルダ

**レスポンス (JSON):**
```json
[
  {
    "id": 1,
    "folder": "example_folder",
    "filename": "example.jpg",
    "timestamp": "2023-11-27T20:40:12Z"
  },
  ...
]
```

---

### POST /image/delete/
**説明:**  
指定された画像を削除し、データベースから画像情報も削除します。  
実装は[handler/delete.go](handler/delete.go)にあります。

**リクエスト (JSON 配列):**
```json
[
  { "id": 1, "folder": "example_folder", "filename": "example.jpg" },
  { "id": 2, "folder": "example_folder", "filename": "sample.png" }
]
```

**レスポンス (JSON 配列):**
```json
[
  { "folder": "example_folder", "filename": "example.jpg" },
  { "folder": "example_folder", "filename": "sample.png" }
]
```

---

### POST /image/move/
**説明:**  
指定された画像のフォルダを移動します。  
移動前のフォルダ、移動後のフォルダおよび対象画像の情報を受け取り、データベースとファイルシステム上の画像パスを更新します。  
詳細は[handler/move-folder.go](handler/move-folder.go)で実装されています。

**リクエスト (JSON):**
```json
{
  "prefolder": "old_folder",
  "postfolder": "new_folder",
  "file": [
    { "id": 1, "filename": "image1.jpg" },
    { "id": 2, "filename": "image2.png" }
  ]
}
```

**レスポンス:**  
成功時はHTTP 200を返します。エラーがあった場合は適切なエラーメッセージが返ります。

---

### POST /image/folder/create/
**説明:**  
新しいフォルダを作成し、ディレクトリおよびデータベースに登録します。  
詳細は[handler/create-folder.go](handler/create-folder.go)で実装されています。

**リクエスト (JSON):**
```json
{
  "folder": "new_folder"
}
```

**レスポンス:**  
作成に成功しても明確なJSONレスポンスは返さず、HTTPステータス200番が返されます。

---

### POST /image/folder/delete/
**説明:**  
指定されたフォルダ内の画像と、データベース上のフォルダ情報を削除します。  
詳細は[handler/delete-folder.go](handler/delete-folder.go)を参照ください。

**リクエスト (JSON 配列):**
```json
[
  { "folder": "folder_to_delete" }
]
```

**レスポンス (JSON 配列):**
```json
[
  { "folder": "folder_to_delete" }
]
```

---

### GET /image/download
**説明:**  
指定された画像ファイルをダウンロードします。  
クエリパラメータでフォルダ、ファイル名、および画像の品質（original/compressed）を指定します。  
実装は[handler/download.go](handler/download.go)にあります。

**リクエストパラメータ (クエリ):**
- **folder**: 画像が保存されているフォルダ名  
- **filename**: ダウンロードする画像ファイル名  
- **quality**: "original" または "compressed"

**レスポンス:**  
対象の画像ファイルが返されます。存在しない場合はエラーメッセージが返ります。

---

## 依存関係

- Go 1.22.3 (toolchain: go1.23.8)
- [github.com/disintegration/imaging](https://pkg.go.dev/github.com/disintegration/imaging)
- [github.com/go-sql-driver/mysql](https://pkg.go.dev/github.com/go-sql-driver/mysql)
- [github.com/google/uuid](https://pkg.go.dev/github.com/google/uuid)
- [github.com/jdeng/goheif](https://pkg.go.dev/github.com/jdeng/goheif)
- [github.com/joho/godotenv](https://pkg.go.dev/github.com/joho/godotenv)

その他、システムレベルで以下が必要です:
- libheif-dev (HEIF形式ファイルの軽量版生成用)

---

## 注意点

- フォルダ名やファイル名は[バリデーション関数](module/validate.go)によりチェックされ、不正な文字列が含まれていると処理が中断されます。  
- 画像の軽量版生成は[module/resize.go](module/resize.go)が担当します。  
- トランザクション処理やファイル操作時のエラーに対しては、ロールバックやエラーログが記録されます。

---

このREADMEは、各エンドポイントの使い方とサーバーの全体的な機能を簡潔に説明しています。詳細な実装や動作確認については、各ソースコードファイルを参照してください。