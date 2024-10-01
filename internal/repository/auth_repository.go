package repository

import (
	"context"
	"database/sql"
	"rest-skeleton/internal/pkg/logger"
)

type AuthRepository struct {
	Db  *sql.DB
	Log *logger.Logger
}

func (r *AuthRepository) HasAuth(ctx context.Context, path string) (bool, error) {
	var hasAuth bool = false

	switch ctx.Err() {
	case context.Canceled:
		return hasAuth, r.Log.Error(ctx, context.Canceled)
	case context.DeadlineExceeded:
		return hasAuth, r.Log.Error(ctx, context.DeadlineExceeded)
	default:
	}

	const q = `
		SELECT 1 
		FROM users 
		JOIN roles_users ON users.id = roles_users.user_id
		JOIN roles ON roles_users.role_id = roles.id
		JOIN access_roles ON roles.id = access_roles.role_id
		JOIN access ON access_roles.access_id = access.id
		WHERE access.path = $1`

	stmt, err := r.Db.PrepareContext(ctx, q)
	if err != nil {
		return hasAuth, r.Log.Error(ctx, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, path).Scan(&hasAuth)
	if err != nil {
		return hasAuth, r.Log.Error(ctx, err)
	}

	return hasAuth, nil
}
