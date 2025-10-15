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
	Order       string    `db:"order_number"`
	Sum         int64     `db:"points"`
	ProcessedAt time.Time `db:"processed_at"`
}

type order struct {
	Number     string    `db:"order_number"`
	Status     string    `db:"status"`
	Accrual    float64   `db:"points"`
	UploadedAt time.Time `db:"uploaded_at"`
}

func accountEntityToDTO(entity account) *dto.Account {
	return &dto.Account{
		Balance:    float64(entity.Balance) / 100,
		TotalSpent: float64(entity.TotalSpent) / 100,
	}
}

func withdrawalFromEntity(entity withdrawal) dto.Withdrawal {
	return dto.Withdrawal{
		Order:       entity.Order,
		Sum:         float64(entity.Sum) / 100,
		ProcessedAt: entity.ProcessedAt.Format(time.RFC3339),
	}
}

func orderFromEntity(entity order) dto.OrderInfo {
	return dto.OrderInfo{
		Number:     entity.Number,
		Status:     entity.Status,
		Accrual:    float64(entity.Accrual) / 100,
		UploadedAt: entity.UploadedAt.Format(time.RFC3339),
	}
}
