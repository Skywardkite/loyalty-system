package handler

import (
	"io"
	"net/http"

	"github.com/Skywardkite/loyalty-system/internal/service"
	"github.com/pkg/errors"
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

	err = h.service.AddOrder(ctx, string(body), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBadRequest):
			http.Error(w, "empty order number", http.StatusBadRequest)
			return
		case errors.Is(err, service.ErrInvalidOrderNumber):
			http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, service.ErrOrderUploadedThisUser):
			w.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, service.ErrOrderUploadedAnotherUser):
			http.Error(w, "order already uploaded by another user", http.StatusConflict)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
