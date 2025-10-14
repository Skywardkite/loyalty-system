package repository

import "github.com/Skywardkite/loyalty-system/internal/handler/dto"

type account struct {
	Balance    int64 `db:"balance"`
	TotalSpent int64 `db:"total_spent"`
}

func entityToDTO(a account) *dto.Account {
	return &dto.Account{
		Balance:    float64(a.Balance) / 100,
		TotalSpent: float64(a.TotalSpent) / 100,
	}
}
