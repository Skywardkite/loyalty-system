package repository

import (
	"time"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
)

type account struct {
	Balance    int64 `db:"balance"`
	TotalSpent int64 `db:"total_spent"`
}

type withdrawal struct {
	Order       string    `db:"order"`
	Sum         int64     `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
}

func accountEntityToDTO(a account) *dto.Account {
	return &dto.Account{
		Balance:    float64(a.Balance) / 100,
		TotalSpent: float64(a.TotalSpent) / 100,
	}
}

func withdrawalFromEntity(w withdrawal) dto.Withdrawal {
	return dto.Withdrawal{
		Order:       w.Order,
		Sum:         float64(w.Sum) / 100,
		ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
	}
}
