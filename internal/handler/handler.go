package handler

import (
	"net/http"

	"github.com/Skywardkite/loyalty-system/internal/repository"
	"github.com/Skywardkite/loyalty-system/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	store   repository.Storage
	logger  *zap.SugaredLogger
	service *service.Service
}

func NewHandler(store repository.Storage, logger *zap.SugaredLogger, s *service.Service) *Handler {
	return &Handler{store: store, logger: logger, service: s}
}

func (h *Handler) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/api/user/register", h.RegisterUser) // регистрация
	r.Post("/api/user/login", h.LoginUser)       // авторизация

	r.Group(func(protected chi.Router) {
		protected.Use(h.AuthMiddleware)

		protected.Get("/api/user/balance", h.GetBalance)                // получение баланса
		protected.Post("/api/user/balance/withdraw", h.WithdrawBalance) // списание баллов
		protected.Get("/api/user/withdrawals", h.GetWithdrawals)        // получение всех списаний бонусов
		protected.Get("/api/user/orders", h.GetOrders)                  // получение загруженных заказов
		protected.Post("/api/user/orders", h.PostOrder)                 // загрузка заказа
	})

	return r
}
