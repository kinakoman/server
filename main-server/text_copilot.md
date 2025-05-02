# サーバーアプリケーション

このプロジェクトは Go 言語で実装された自宅サーバー用のアプリケーションです。  
リバースプロキシとして動作し、各種バックエンドサーバー（ホームページサーバー、image-server）へリクエストを転送します。  
また、ログイン認証、CSRFトークン管理、アクセスログのデータベース記録、MySQL停止時のバックアップ機能などを備えています。

## API エンドポイントと機能

### 全般
- **全リクエスト**  
  すべてのエンドポイントは、認証チェックが実施されます。  
  ブラウザからの未認証リクエストは、`/login` のログイン画面（[template/login.html](template/login.html)）が返されます。

### 認証関連
- **GET /login**  
  ログイン画面を表示します。詳細は [`auth.LoginHandler`](auth/auth-handler.go) を参照してください。

- **POST /login**  
  フォームから送信された `username` と `password` を元に認証を行い、  
  認証が成功した場合はセッションIDが発行され、クッキーに設定されます。  
  認証失敗の場合はログイン画面が再表示されます。

- **GET /csrf-token**  
  セッションに紐付いた CSRF トークンを取得します。  
  このエンドポイントは CSRF チェックをスキップしています。  
  詳細は [`auth.GetCsrfTokenHandler`](auth/auth-handler.go) を参照してください。

- **/logout**  
  セッションを破棄し、ログアウトを実施します。  
  詳細は [`auth.LogoutHandler`](auth/auth-handler.go) を参照してください。

### プロキシ
- **/ (root)**  
  認証が完了している場合、[`module.SetMiddleware`](module/middleware.go) により、  
  ホームページサーバーへのリバースプロキシ（[`proxy.InitHomepageProxy`](proxy/homepage-proxy.go)）が適用されます。

- **/image/**  
  画像サーバーへのリバースプロキシ機能を提供します。  
  リクエストヘッダーに秘密鍵情報を付与し、200 番以外のレスポンスは 502 番に書き換える仕組みが（必要に応じて）用意されています。  
  詳細は [`proxy.InitImageProxy`](proxy/image-proxy.go) を参照してください。

### サーバーステータス
- **POST /status/**  
  サーバーモニターからのアクセス用で、Basic 認証（環境変数 `LOCAL_AUTH_USER`、`LOCAL_AUTH_PASSWORD` を使用）で保護されています。  
  アクセスログは記録されず、認証成功時のみ HTTP 200 を返します。  
  詳細は [`handler.StatusHandler`](handler/status.go) を参照してください。

## セットアップ

1. **Go のインストール**  
   Go 言語のバージョン 1.23 以上が必要です。

2. **依存パッケージのインストール**  
   プロジェクトルートで以下を実行します:
   
    ```sh
    go mod download
    ```

3. **環境変数の設定**  
   プロジェクトのルートまたは `../.env` ファイルに以下の変数を定義してください（例）:
   
    ```dotenv
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_HOST=localhost
    DB_PORT=3306
    DB_NAME=your_db_name
    DB_TABLE=server_log

    AUTH_USER_TABLE=auth_user
    AUTH_SESSION_TABLE=auth_session

    COOKIE_SESSION_NAME=your_session_cookie
    PROXY_SECRET_HEADER=your_proxy_header
    PROXY_SECRET_KEY=your_secret_key

    MAIN_SERVER_HOST=0.0.0.0
    MAIN_SERVER_PORT=8080
    MAIN_HOMEPAGE_HOST=homepage_host
    MAIN_HOMEPAGE_PORT=homepage_port

    SERVER_SCHEME=http

    IMAGE_SERVER_HOST=image_server_host
    IMAGE_SERVER_PORT=image_server_port

    LOCAL_AUTH_USER=local_admin
    LOCAL_AUTH_PASSWORD=local_secret
    ```

4. **MySQL データベースの作成**  
   - `server_log` テーブル：アクセスログ、バックアップログ用  
   - `auth_user` テーブル：ユーザー情報（username, password）  
   - `auth_session` テーブル：セッション情報（session_id, csrf_token, user_id, expires_at）  

   各テーブルのカラム構成は [README.md の記述](README.md) と同様です。

## 起動方法

1. **サーバーの起動**  
   プロジェクトの `main.go` がエントリーポイントです。  
   以下のコマンドでサーバーを起動します:

    ```sh
    go run main.go
    ```

2. **動作確認**  
   - ブラウザで `http://<MAIN_SERVER_HOST>:<MAIN_SERVER_PORT>/login` にアクセスし、ログイン画面が表示されることを確認します。
   - 認証後、`/` ではホームページサーバーへのプロキシが動作し、`/image/` では画像サーバーへのプロキシが動作します。
   - サーバー監視用のステータス更新は、`POST /status/` に対してBasic認証（`LOCAL_AUTH_USER` / `LOCAL_AUTH_PASSWORD`）でリクエストを送ることで動作を確認できます。


## 補足

- **バックアップ機能**  
  MySQL の停止時は、アクセスログはバックアップ用に [`module/backup.go`](module/backup.go) で管理されます。一定時間ごとに MySQL の状態を確認し、起動時にバックアップログをデータベースへ書き込みます。

- **ミドルウェア**  
  すべてのリクエストには、[`module/middleware.go`](module/middleware.go) で定義されたログ記録、認証、CSRFチェックのミドルウェアが適用されます。

- **環境変数**  
  環境変数は `.env` ファイル等で適宜管理してください。プロジェクト起動時に [`github.com/joho/godotenv`](go.mod) が読み込みます。

この README により、プロジェクトの全体像や起動手順、API エンドポイントの仕様を把握いただけます。
