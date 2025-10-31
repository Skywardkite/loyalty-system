package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	"github.com/Skywardkite/loyalty-system/internal/service"
	"github.com/pkg/errors"
)

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID, err := h.service.Login(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBadRequest):
			http.Error(w, "login and password required", http.StatusBadRequest)
			return
		case errors.Is(err, service.ErrUnauthorized):
			http.Error(w, "invalid login or password", http.StatusUnauthorized)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	h.getToken(userID, w)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getToken(userID int64, w http.ResponseWriter) {
	// Генерация токена
	token, err := service.GenerateAuthToken(userID)
	if err != nil {
		h.logger.Errorw("failed to generate token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Устанавливаем cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
