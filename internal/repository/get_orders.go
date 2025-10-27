package repository

import (
	"context"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	slices "github.com/Skywardkite/loyalty-system/internal/utils"
)

func (r *Repository) GetOrdersByUser(ctx context.Context, userID int64) ([]dto.OrderInfo, error) {
	var rows []order

	query := `
		SELECT order_number, status, points, uploaded_at
		FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at DESC
	`
	if err := r.db.SelectContext(ctx, &rows, query, userID); err != nil {
		return nil, err
	}

	return slices.Map(rows, orderInfoFromEntity), nil
}
