package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	repErr "github.com/Skywardkite/loyalty-system/internal/repository/error"
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

	if req.Order == "" {
		http.Error(w, "invalid order", http.StatusUnprocessableEntity)
		return
	}

	if req.Sum <= 0 {
		http.Error(w, "invalid sum", http.StatusUnprocessableEntity)
		return
	}

	sum := int64(req.Sum * 100)

	if err := h.store.AddWithdrawal(ctx, userID, sum, req.Order); err != nil {
		switch {
		case errors.Is(err, repErr.ErrInsufficientFunds):
			http.Error(w, "insufficient funds", http.StatusPaymentRequired)
			return
		default:
			h.logger.Errorw("failed to create withdrawal", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
