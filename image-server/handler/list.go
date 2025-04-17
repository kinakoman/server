package handler

import (
	"net/http"
)

type ListHandler struct{}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
