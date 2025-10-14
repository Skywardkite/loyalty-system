package repository

import "context"

func (r *Repository) UserExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE login=$1)", login).Scan(&exists)
	return exists, err
}
