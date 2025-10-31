package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	repErr "github.com/Skywardkite/loyalty-system/internal/repository/error"
	"github.com/Skywardkite/loyalty-system/internal/service"
)

func (h *Handler) WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.ParamsWithdraw

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.service.WriteOffPoints(ctx, &req, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnprocessableEntity):
			http.Error(w, "invalid data in request", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, repErr.ErrInsufficientFunds):
			http.Error(w, "insufficient funds", http.StatusPaymentRequired)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
