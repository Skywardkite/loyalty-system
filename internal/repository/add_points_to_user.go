package repository

import (
	"context"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
)

func (r *Repository) AddPointsToUser(ctx context.Context, order dto.Order) error {
	entity := orderToEntity(order)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = r.db.ExecContext(ctx, `
		UPDATE accounts SET balance = balance + $1 WHERE user_id = $2`,
		entity.Accrual, entity.UserID)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE orders SET status = $1, points = $2 WHERE order_number = $3`,
		entity.Status, entity.Accrual, entity.Number)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
