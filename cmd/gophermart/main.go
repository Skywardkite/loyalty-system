package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Skywardkite/loyalty-system/internal/accrual"
	"github.com/Skywardkite/loyalty-system/internal/config"
	"github.com/Skywardkite/loyalty-system/internal/handler"
	"github.com/Skywardkite/loyalty-system/internal/repository"
	"github.com/Skywardkite/loyalty-system/internal/service"
	"github.com/Skywardkite/loyalty-system/pkg/logger"
)

func main() {
	if err := logger.Initialize(); err != nil {
		log.Fatal("Error to initialize logger:", err)
	}
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.ParseFlags()
	if err != nil {
		logger.Sugar.Fatalw("Error to parse flags", "error", err)
	}

	store, err := repository.New(cfg.DatabaseURI)
	if err != nil {
		logger.Sugar.Fatalw("Failed to connect to database", "error", err)
	}

	defer func() {
		if err := store.Close(); err != nil {
			logger.Sugar.Errorw("Failed to close database", "error", err)
		}
	}()

	accrualClient := accrual.NewClient(cfg.AccrualSystemAddress, 5*time.Second)
	s := service.New(logger.Sugar, store, accrualClient)
	h := handler.NewHandler(store, logger.Sugar, s)

	// Запускаем воркер для постоянной обработки заказов через accrual
	go s.StartAccrualWorker(ctx)

	srv := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: h.RegisterRoutes(),
	}

	logger.Sugar.Infow("Starting server", "addr", cfg.RunAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Sugar.Fatalw("Server error", "error", err)
	}
}
