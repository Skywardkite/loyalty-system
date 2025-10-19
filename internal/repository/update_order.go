package repository

import (
	"context"
)

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderNumber, status string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE orders SET status = $1 WHERE order_number = $2`,
		status, orderNumber)

	return err
}
