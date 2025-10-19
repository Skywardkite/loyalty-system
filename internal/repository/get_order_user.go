package repository

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	repErr "github.com/Skywardkite/loyalty-system/internal/repository/error"
)

func (r *Repository) GetOrderUserByNumber(ctx context.Context, order string) (userID int64, err error) {
	err = r.db.QueryRowContext(ctx,
		"SELECT user_id FROM orders WHERE order_number=$1",
		order).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, repErr.ErrOrderNotExsit
	}

	return
}
