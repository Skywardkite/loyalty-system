package repository

import (
	"context"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	slices "github.com/Skywardkite/loyalty-system/internal/utils"
)

func (r *Repository) GetUnprocessedOrders(ctx context.Context) ([]dto.Order, error) {
	var rows []order

	query := `SELECT order_number, user_id, status 
		FROM orders 
		WHERE status NOT IN ('PROCESSED', 'INVALID')`

	err := r.db.SelectContext(ctx, &rows, query)

	return slices.Map(rows, orderFromEntity), err
}
