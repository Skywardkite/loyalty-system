package repository

import "context"

func (r *Repository) CreateUser(ctx context.Context, login, password string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id",
		login, password).Scan(&id)
	if err != nil {
		return 0, err
	}

	// создаём аккаунт для хранения баллов
	_, err = r.db.ExecContext(ctx, "INSERT INTO accounts (user_id) VALUES ($1)", id)
	if err != nil {
		return id, err
	}

	return id, nil
}
