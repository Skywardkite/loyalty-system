package repository

import (
	"context"
	"fmt"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	UserExists(ctx context.Context, login string) (bool, error)
	GetUserByLogin(ctx context.Context, login string) (id int64, hashedPassword string, err error)
	GetBalance(ctx context.Context, userID int64) (*dto.Account, error)
	AddWithdrawal(ctx context.Context, userID, sum int64, orderID string) error
	GetWithdrawalsByUser(ctx context.Context, userID int64) ([]dto.Withdrawal, error)
	GetOrdersByUser(ctx context.Context, userID int64) ([]dto.OrderInfo, error)
	AddOrder(ctx context.Context, userID int64, order string) error
	GetOrderUserByNumber(ctx context.Context, order string) (userID int64, err error)
	GetUnprocessedOrders(ctx context.Context) ([]dto.Order, error)
	UpdateOrderStatus(ctx context.Context, orderNumber, status string) error
	AddPointsToUser(ctx context.Context, order dto.Order) error
}

type Repository struct {
	db *sqlx.DB
}

func New(uri string) (*Repository, error) {
	db, err := sqlx.Connect("pgx", uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = applyMigrations(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	if r.db == nil {
		return nil
	}
	return r.db.Close()
}

func applyMigrations(uri string) error {
	m, err := migrate.New(
		"file://migrations",
		uri,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Применяем все новые миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
