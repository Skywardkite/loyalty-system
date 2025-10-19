package handler

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Skywardkite/loyalty-system/internal/service"
)

func (h *Handler) PostOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	orderNum := string(body)

	number := strings.TrimSpace(orderNum)
	if number == "" {
		http.Error(w, "empty order number", http.StatusBadRequest)
		return
	}

	err = h.service.AddOrder(ctx, userID, orderNum)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderNotExsit):
			http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, service.ErrUploadedAnotherUser):
			http.Error(w, "order already uploaded by another user", http.StatusConflict)
			return
		case errors.Is(err, service.ErrUploadedUser):
			w.WriteHeader(http.StatusOK)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
