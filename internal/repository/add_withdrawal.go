package repository

import (
	"context"

	repErr "github.com/Skywardkite/loyalty-system/internal/repository/error"
)

func (r *Repository) AddWithdrawal(ctx context.Context, userID, sum int64, orderID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var balance float32

	// Блокируем баланс пользователя до конца транзакции
	err = tx.GetContext(ctx, &balance, `
		SELECT balance 
		FROM accounts 
		WHERE user_id = $1 
		FOR UPDATE
	`, userID)
	if err != nil {
		return err
	}

	if int64(balance) < sum {
		return repErr.ErrInsufficientFunds
	}

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

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
