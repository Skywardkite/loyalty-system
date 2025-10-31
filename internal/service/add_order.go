package service

import (
	"context"
	"strings"

	repErr "github.com/Skywardkite/loyalty-system/internal/repository/error"
	"github.com/pkg/errors"
)

var (
	ErrInvalidOrderNumber       = errors.New("Невалидный номер заказа")
	ErrOrderUploadedThisUser    = errors.New("Заказ уже добавлен этим пользователем")
	ErrOrderUploadedAnotherUser = errors.New("Заказ уже загружен другим пользователем")
)

func (s *Service) AddOrder(ctx context.Context, number string, userID int64) error {
	orderNumber := strings.TrimSpace(number)
	if number == "" {
		return ErrBadRequest
	}

	if !validateOrder(orderNumber) {
		return ErrInvalidOrderNumber
	}

	orderUserID, err := s.store.GetOrderUserByNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, repErr.ErrOrderNotExsit) {
		s.logger.Errorw("failed to check order", "error", err)
		return err
	}

	if orderUserID != 0 {
		if orderUserID == userID {
			return ErrOrderUploadedThisUser
		} else {
			return ErrOrderUploadedAnotherUser
		}
	}

	err = s.store.AddOrder(ctx, userID, orderNumber)
	if err != nil {
		s.logger.Errorw("failed to add order", "error", err)
		return err
	}

	return nil
}
