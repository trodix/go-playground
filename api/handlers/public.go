package handlers

import (
	"encoding/json"
	"net/http"
)

type PublicHandler struct {
}

func NewPublicHandler() *PublicHandler {
	return &PublicHandler{}
}

func (h *PublicHandler) Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("World")
}
