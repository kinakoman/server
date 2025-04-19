package handler

import (
	"net/http"
)

// "/"アクセスのハンドラ
type RootHandler struct{}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "This is KINAKO Server.\n")
	// fmt.Fprintf(w, "/homepage/ : You can access my homepage.\n")
	w.Write([]byte("<!DOCTYPE html>\n" +
		"<html lang=\"en\">\n" +
		"<head><meta charset=\"UTF-8\"><title>Homepage</title></head>\n" +
		"<body>\n" +
		"<h2>Welcome to the Homepage!</h2>\n" +
		"<form method=\"POST\" action=\"/logout\">\n" +
		"<button type=\"submit\">Logout</button>\n" +
		"</form>\n" +
		"<form action=\"http://localhost:8080/image/upload/\" method=\"post\" enctype=\"multipart/form-data\">\n" +
		"<label>画像アップロード:</label><br>\n" +
		"<input type=\"file\" name=\"images\" multiple><br><br>\n" +
		"<label>保存先フォルダ名:</label><br>\n" +
		"<input type=\"text\" name=\"folder\"><br><br>\n" +
		"<input type=\"submit\" value=\"アップロード\">\n" +
		"</form>\n" +
		"<h1>Delete Files</h1>\n" +
		"<div id=\"formContainer\">\n" +
		"    <div class=\"entry\">\n" +
		"        <label>Folder: <input type=\"text\" class=\"folder\" required></label>\n" +
		"        <label>Filename: <input type=\"text\" class=\"filename\" required></label>\n" +
		"    </div>\n" +
		"</div>\n" +
		"<button type=\"button\" id=\"addBtn\">＋ 追加</button>\n" +
		"<button type=\"button\" id=\"submitBtn\">送信</button>\n" +
		"<h1>フォルダ削除</h1>\n" +
		"<form id=\"deleteFolderForm\">\n" +
		"  <label>削除対象フォルダ名:</label>\n" +
		"  <input type=\"text\" id=\"deleteFolder\" name=\"folder\" required>\n" +
		"  <button type=\"submit\">削除</button>\n" +
		"</form>\n" +
		"<script>\n" +
		"document.getElementById('addBtn').addEventListener('click', function() {\n" +
		"    const container = document.getElementById('formContainer');\n" +
		"    const entry = document.createElement('div');\n" +
		"    entry.className = 'entry';\n" +
		"    entry.innerHTML = '<label>Folder: <input type=\\\"text\\\" class=\\\"folder\\\" required></label>' +\n" +
		"                      '<label>Filename: <input type=\\\"text\\\" class=\\\"filename\\\" required></label>';\n" +
		"    container.appendChild(entry);\n" +
		"});\n" +
		"document.getElementById('submitBtn').addEventListener('click', function() {\n" +
		"    const folders = document.querySelectorAll('.folder');\n" +
		"    const filenames = document.querySelectorAll('.filename');\n" +
		"    const data = [];\n" +
		"    for (let i = 0; i < folders.length; i++) {\n" +
		"        const folder = folders[i].value.trim();\n" +
		"        const filename = filenames[i].value.trim();\n" +
		"        if (folder && filename) {\n" +
		"            data.push({ folder: folder, filename: filename });\n" +
		"        }\n" +
		"    }\n" +
		"    if (data.length === 0) return;\n" +
		"    fetch('/image/delete/', {\n" +
		"        method: 'POST',\n" +
		"        headers: { 'Content-Type': 'application/json' },\n" +
		"        body: JSON.stringify(data)\n" +
		"    });\n" +
		"});\n" +
		"document.getElementById('deleteFolderForm').addEventListener('submit', function(e) {\n" +
		"  e.preventDefault();\n" +
		"  const folder = document.getElementById('deleteFolder').value.trim();\n" +
		"  if (!folder) return;\n" +
		"  fetch('/image/folder/delete/', {\n" +
		"    method: 'POST',\n" +
		"    headers: { 'Content-Type': 'application/json' },\n" +
		"    body: JSON.stringify({ folder: folder })\n" +
		"  });\n" +
		"});\n" +
		"</script>\n" +
		"</body>\n" +
		"</html>"))
}
