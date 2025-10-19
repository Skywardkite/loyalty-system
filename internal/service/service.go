package service

import (
	"github.com/Skywardkite/loyalty-system/internal/accrual"
	"github.com/Skywardkite/loyalty-system/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	logger        *zap.SugaredLogger
	store         repository.Storage
	accrualClient *accrual.Client
}

func New(logger *zap.SugaredLogger, s repository.Storage,
	accrualClient *accrual.Client) *Service {
	return &Service{logger: logger, store: s, accrualClient: accrualClient}
}
