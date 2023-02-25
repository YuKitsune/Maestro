package handlers

import (
	"github.com/yukitsune/maestro/pkg/api/responses"
	"net/http"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	responses.NotFound(w, "route not found")
}
