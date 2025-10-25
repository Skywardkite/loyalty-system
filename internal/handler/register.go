package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	"github.com/Skywardkite/loyalty-system/internal/service"
	"github.com/pkg/errors"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID, err := h.service.Registration(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBadRequest):
			http.Error(w, "login and password required", http.StatusBadRequest)
			return
		case errors.Is(err, service.ErrUserAlredyExists):
			http.Error(w, "login already exists", http.StatusConflict)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	h.getToken(userID, w)

	w.WriteHeader(http.StatusOK)
}
