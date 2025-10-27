package repository

import "github.com/pkg/errors"

var (
	ErrOrderNotExsit     = errors.New("Заказа в системе нет")
	ErrInsufficientFunds = errors.New("Недостаточно средств на балансе")
)
