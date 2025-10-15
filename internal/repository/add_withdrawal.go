package repository

import (
	"context"
)

func (r *Repository) AddWithdrawal(ctx context.Context, userID, sum int64, orderID string) error {
	// Транзакция, чтобы одинаковые данные по списаниям были
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = r.db.ExecContext(ctx, `
		UPDATE accounts
		SET balance = balance - $1, total_spent = total_spent + $2
		WHERE user_id = $3`,
		sum, sum, userID)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO withdrawals 
			(user_id, order_number, points) 
		VALUES ($1, $2, $3)`,
		userID, orderID, sum)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
