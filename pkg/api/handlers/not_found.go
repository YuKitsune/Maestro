package handlers

import (
	"net/http"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	NotFound(w, "route not found")
}
