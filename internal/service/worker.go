package service

import (
	"context"
	"net/http"
	"time"

	"github.com/Skywardkite/loyalty-system/internal/constants"
)

func (s *Service) StartAccrualWorker(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	s.logger.Infow("accrual worker started")

	for {
		select {
		case <-ctx.Done():
			s.logger.Infow("accrual worker stopped")
			return
		case <-ticker.C:
			s.processOrders(ctx)
		}
	}
}

func (s *Service) processOrders(ctx context.Context) {
	orders, err := s.store.GetUnprocessedOrders(ctx)
	if err != nil {
		s.logger.Errorw("failed to get unprocessed orders", "error", err)
		return
	}

	for _, order := range orders {
		info, code, retryAfter, err := s.accrualClient.GetOrderAccrual(ctx, order.Number)
		if err != nil {
			s.logger.Errorw("failed to fetch order from accrual", "order", order.Number, "error", err)
			continue
		}

		switch code {
		case http.StatusTooManyRequests:
			s.logger.Warnw("accrual rate limit hit", "order", order.Number)
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		case http.StatusNoContent:
			continue
		case http.StatusOK:
			if info.Status == order.Status {
				continue
			}

			order.Status = info.Status

			if info.Status == constants.Processed {
				err = s.store.AddPointsToUser(ctx, order)
				if err != nil {
					s.logger.Errorw("failed to add points to user", "order", order.Number, "error", err)
				}

				continue
			}

			err = s.store.UpdateOrderStatus(ctx, order.Number, constants.GetStatusOrder(order.Status))
			if err != nil {
				s.logger.Errorw("failed to update order status", "order", order.Number, "error", err)
			}
		default:
			s.logger.Warnw("unexpected accrual response", "order", order.Number, "code", code)
		}
	}
}
