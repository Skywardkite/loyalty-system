package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	"github.com/Skywardkite/loyalty-system/internal/service"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password required", http.StatusBadRequest)
		return
	}

	userID, hashedPassword, err := h.store.GetUserByLogin(ctx, req.Login)
	if err != nil {
		h.logger.Errorw("failed to get user by login", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Сравниваем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "invalid login or password", http.StatusUnauthorized)
		return
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
