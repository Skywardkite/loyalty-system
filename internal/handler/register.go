package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
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

	exists, err := h.store.UserExists(ctx, req.Login)
	if err != nil {
		h.logger.Errorw("failed to check user existence", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "login already exists", http.StatusConflict)
		return
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Errorw("failed to hash password", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Создаём пользователя
	userID, err := h.store.CreateUser(ctx, req.Login, string(hashedPassword))
	if err != nil {
		h.logger.Errorw("failed to create user", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.getToken(userID, w)

	w.WriteHeader(http.StatusOK)
}
