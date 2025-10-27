package handler

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
)

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.store.GetWithdrawalsByUser(ctx, userID)
	if err != nil {
		h.logger.Errorw("failed to get withdrawals", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")

	gz := gzip.NewWriter(w)
	defer gz.Close()

	if err := json.NewEncoder(gz).Encode(withdrawals); err != nil {
		h.logger.Errorw("failed to encode withdrawals", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
