package repository

import "context"

func (r *Repository) AddOrder(ctx context.Context, userID int64, order string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO orders (order_number, user_id) VALUES ($1, $2)", order, userID)

	return err
}
