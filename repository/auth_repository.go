package repository

import (
	"context"
	"database/sql"
	"errors"
	"rest-skeleton/pkg/logger"
)

type AuthRepository struct {
	Db  *sql.DB
	Log *logger.Logger
}

func (r *AuthRepository) HasAuth(ctx context.Context, path string) (bool, error) {
	var hasAuth bool = false

	switch ctx.Err() {
	case context.Canceled:
		err := errors.New("request is canceled")
		r.Log.Error.Println(err)
		return hasAuth, err
	case context.DeadlineExceeded:
		err := errors.New("deadline is exceeded")
		r.Log.Error.Println(err)
		return hasAuth, err
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
		r.Log.Error.Println("error get users", err)
		return hasAuth, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, path).Scan(&hasAuth)
	if err != nil {
		r.Log.Error.Println("error get users", err)
		return hasAuth, err
	}

	return hasAuth, nil
}
