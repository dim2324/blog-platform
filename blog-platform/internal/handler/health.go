package handler

import (
	"net/http"

	"blog-platform/pkg/httpjson"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
