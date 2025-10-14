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

	r.Post("/api/user/register", h.RegisterUser)
	r.Post("/api/user/login", h.LoginUser)

	r.Group(func(chi.Router) {
		r.Get("/api/user/balance", h.GetBalance)
	})

	return r
}