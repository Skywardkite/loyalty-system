package repository

import "context"

func (r *Repository) GetUserByLogin(ctx context.Context, login string) (id int64, hashedPassword string, err error) {
	err = r.db.QueryRowContext(ctx,
		"SELECT id, password FROM users WHERE login=$1",
		login).Scan(&id, &hashedPassword)
	return
}
