package handler

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	repErr "github.com/Skywardkite/loyalty-system/internal/repository/error"
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

	isDone := h.addOrder(ctx, w, userID, orderNum)
	if isDone {
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) addOrder(ctx context.Context, w http.ResponseWriter, userID int64, orderNumber string) bool {
	if !service.ValidateOrder(orderNumber) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return true
	}

	orderUserID, err := h.store.GetOrderUserByNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, repErr.ErrOrderNotExsit) {
		h.logger.Errorw("failed to check order", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return true
	}

	if orderUserID != 0 {
		if orderUserID == userID {
			w.WriteHeader(http.StatusOK)
			return true
		} else {
			http.Error(w, "order already uploaded by another user", http.StatusConflict)
			return true
		}
	}

	err = h.store.AddOrder(ctx, userID, orderNumber)
	if err != nil {
		h.logger.Errorw("failed to add order", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return true
	}

	return false
}
