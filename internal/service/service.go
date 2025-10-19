package service

import (
	"github.com/Skywardkite/loyalty-system/internal/config"
	"github.com/Skywardkite/loyalty-system/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	Cfg    *config.Config
	logger *zap.SugaredLogger
	store  repository.Storage
}

func New(cfg *config.Config, logger *zap.SugaredLogger, s repository.Storage) *Service {
	return &Service{Cfg: cfg, logger: logger, store: s}
}
