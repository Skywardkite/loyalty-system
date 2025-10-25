package service

import (
	"context"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	repErr "github.com/Skywardkite/loyalty-system/internal/repository/error"
	"github.com/pkg/errors"
)

var (
	ErrUnprocessableEntity = errors.New("Невалидные данные в запросе")
)

func (s *Service) WriteOffPoints(ctx context.Context, req *dto.ParamsWithdraw, userID int64) error {
	if req.Order == "" {
		return ErrUnprocessableEntity
	}

	if req.Sum <= 0 {
		return ErrUnprocessableEntity
	}

	sum := int64(req.Sum * 100)

	err := s.store.AddWithdrawal(ctx, userID, sum, req.Order)
	if !errors.Is(err, repErr.ErrInsufficientFunds) {
		s.logger.Errorw("failed to create withdrawal", "error", err)
	}

	return err
}
