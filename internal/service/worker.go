package service

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Skywardkite/loyalty-system/internal/accrual"
	"github.com/Skywardkite/loyalty-system/internal/constants"
	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
)

func (s *Service) StartAccrualWorker(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
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

	var wg sync.WaitGroup
	pauseChan := make(chan time.Duration, 1) // канал для паузы при 429

	for _, order := range orders {
		wg.Add(1)

		go func() {
			defer wg.Done()

			select {
			case pause := <-pauseChan:
				s.logger.Warnw("pausing worker due to 429", "wait", pause.Seconds())
				time.Sleep(pause)
			default:
				// можем идти дальше, пока StatusTooManyRequests не было
			}

			info, code, retryAfter, err := s.accrualClient.GetOrderAccrual(ctx, order.Number)
			if err != nil {
				s.logger.Errorw("failed to fetch order from accrual", "order", order.Number, "error", err)
				return
			}

			switch code {
			case http.StatusTooManyRequests:
				s.logger.Warnw("accrual rate limit hit", "order", order.Number)
				select {
				case pauseChan <- time.Duration(retryAfter) * time.Second:
				default:
					// пропускаем, если канал уже содержит паузу
				}
				return
			case http.StatusNoContent:
				return
			case http.StatusOK:
				s.updateOrder(ctx, info, &order)
			default:
				s.logger.Warnw("unexpected accrual response", "order", order.Number, "code", code)
			}
		}()
	}

	wg.Wait()
}

func (s *Service) updateOrder(ctx context.Context, info *accrual.OrderInfo, order *dto.Order) {
	if info.Status == order.Status {
		return
	}

	order.Status = info.Status

	if info.Status == constants.Processed {
		if info.Accrual != nil && *info.Accrual > 0 {
			err := s.store.AddPointsToUser(ctx, *order)
			if err != nil {
				s.logger.Errorw("failed to add points to user", "order", order.Number, "error", err)
			}

			return
		}
	}

	err := s.store.UpdateOrderStatus(ctx, order.Number, constants.GetStatusOrder(order.Status))
	if err != nil {
		s.logger.Errorw("failed to update order status", "order", order.Number, "error", err)
	}
}
