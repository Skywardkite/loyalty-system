package repository

import (
	"context"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
)

func (r *Repository) GetBalance(ctx context.Context, userID int64) (*dto.Account, error) {
	var a account
	err := r.db.GetContext(ctx, &a, "SELECT balance, total_spent FROM accounts WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}

	return accountEntityToDTO(a), nil
}
