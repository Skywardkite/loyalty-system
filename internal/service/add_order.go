package service

import (
	"context"
	"errors"

	repErr "github.com/Skywardkite/loyalty-system/internal/repository/error"
)

var (
	ErrOrderNotExsit       = errors.New("invalid order number")
	ErrUploadedAnotherUser = errors.New("uploaded by another user")
	ErrUploadedUser        = errors.New("uploaded by user")
)

func (s *Service) AddOrder(ctx context.Context, userID int64, orderNumber string) error {
	err := validateOrder(orderNumber)
	if err != nil {
		return err
	}

	orderUserID, err := s.store.GetOrderUserByNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, repErr.ErrOrderNotExsit) {
		s.logger.Errorw("failed to check order", "error", err)
		return err
	}

	if orderUserID != 0 {
		if orderUserID == userID {
			return ErrUploadedUser
		} else {
			return ErrUploadedAnotherUser
		}
	}

	err = s.store.AddOrder(ctx, userID, orderNumber)
	if err != nil {
		s.logger.Errorw("failed to get balance", "error", err)
		return err
	}

	return nil
}
