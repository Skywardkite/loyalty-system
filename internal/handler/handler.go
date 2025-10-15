package handler

import (
	"net/http"

	"github.com/Skywardkite/loyalty-system/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	store  repository.Storage
	logger *zap.SugaredLogger
}

func NewHandler(store repository.Storage, logger *zap.SugaredLogger) *Handler {
	return &Handler{store: store, logger: logger}
}

func (h *Handler) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/api/user/register", h.RegisterUser) // регистрация
	r.Post("/api/user/login", h.LoginUser)       // авторизация

	r.Group(func(chi.Router) {
		r.Get("/api/user/balance", h.GetBalance)                // получение баланса
		r.Post("/api/user/balance/withdraw", h.WithdrawBalance) // списание баллов
		r.Get("/api/user/withdrawals", h.GetWithdrawals)        // получение всех списаний бонусов
	})

	return r
}
