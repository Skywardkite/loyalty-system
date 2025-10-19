package handler

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

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

	order := h.validateOrder(w, body)
	if order != "" {
		return
	}

	orderUserID, err := h.store.GetOrderUserByNumber(ctx, order)
	if err != nil && !errors.Is(err, repErr.ErrOrderNotExsit) {
		h.logger.Errorw("failed to check order", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if orderUserID != 0 {
		if orderUserID == userID {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			http.Error(w, "order already uploaded by another user", http.StatusConflict)
			return
		}
	}

	err = h.store.AddOrder(ctx, userID, order)
	if err != nil {
		h.logger.Errorw("failed to get balance", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	//TODO: горутина для взаимодействия с сервисом accrual
}

func (h *Handler) validateOrder(w http.ResponseWriter, body []byte) string {
	orderNum := strings.TrimSpace(string(body))
	if orderNum == "" {
		http.Error(w, "empty order number", http.StatusBadRequest)
		return ""
	}

	// Проверяем, что только цифры
	if !regexp.MustCompile(`^\d+$`).MatchString(orderNum) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return ""
	}

	// Проверка алгоритмом Луна
	if !service.IsValidLuhn(orderNum) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return ""
	}

	return orderNum
}