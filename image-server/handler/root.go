package handler

import (
	"fmt"
	"net/http"
)

type RootHandler struct{}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is KINAKOSAMA Server.\n")
}
