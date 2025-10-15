package repository

import (
	"context"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	slices "github.com/Skywardkite/loyalty-system/internal/utils"
)

func (r *Repository) GetWithdrawalsByUser(ctx context.Context, userID int64) ([]dto.Withdrawal, error) {
	var rows []withdrawal

	query := `
		SELECT order_number, sum, processed_at
		FROM withdrawals
		WHERE user_id = $1
		ORDER BY processed_at DESC
	`
	if err := r.db.SelectContext(ctx, &rows, query, userID); err != nil {
		return nil, err
	}

	return slices.Map(rows, withdrawalFromEntity), nil
}
