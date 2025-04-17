package handler

import (
	"net/http"
)

// "/"アクセスのハンドラ
type RootHandler struct{}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "This is KINAKO Server.\n")
	// fmt.Fprintf(w, "/homepage/ : You can access my homepage.\n")
	w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Homepage</title>
</head>
<body>
  <h2>Welcome to the Homepage!</h2>
  <form method="POST" action="/logout">
    <button type="submit">Logout</button>
  </form>
  <form action="http://localhost:8080/image/upload/" method="post" enctype="multipart/form-data">
  <label>画像アップロード:</label><br>
  <input type="file" name="images" multiple><br><br>

  <label>保存先フォルダ名:</label><br>
  <input type="text" name="folder"><br><br>

  <input type="submit" value="アップロード">
</form>
</body>
</html>`))
}
