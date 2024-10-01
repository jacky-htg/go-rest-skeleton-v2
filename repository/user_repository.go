package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"rest-skeleton/model"
	"rest-skeleton/pkg/logger"
	"rest-skeleton/pkg/myctx"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserRepository struct {
	Db         *sql.DB
	Log        *logger.Logger
	UserEntity model.User
}

func (u *UserRepository) Find(ctx context.Context) error {
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "FindUserRepository")
	defer span.End()

	switch ctx.Err() {
	case context.Canceled:
		return u.Log.Error(ctx, context.Canceled)
	case context.DeadlineExceeded:
		return u.Log.Error(ctx, context.DeadlineExceeded)
	default:
	}

	const q = `SELECT ids, name, email, password FROM users WHERE id=$1 AND deleted_at IS NULL`
	span.SetAttributes(attribute.String("db.query", q))
	span.SetAttributes(attribute.Int64("db.id", u.UserEntity.ID))

	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		return u.Log.Error(ctx, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.UserEntity.ID).Scan(&u.UserEntity.ID, &u.UserEntity.Name, &u.UserEntity.Email, &u.UserEntity.Password)
	if err != nil {
		return u.Log.Error(ctx, err)
	}
	return nil
}

func (u *UserRepository) Save(ctx context.Context) error {
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "SaveUserRepository")
	defer span.End()

	switch ctx.Err() {
	case context.Canceled:
		return u.Log.Error(ctx, context.Canceled)
	case context.DeadlineExceeded:
		return u.Log.Error(ctx, context.DeadlineExceeded)
	default:
	}

	const q = `INSERT INTO users (name, password, email, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	span.SetAttributes(attribute.String("db.query", q))

	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		return u.Log.Error(ctx, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		u.UserEntity.Name,
		u.UserEntity.Password,
		u.UserEntity.Email,
		ctx.Value(myctx.Key("user_id")).(int64),
	).Scan(&u.UserEntity.ID)
	if err != nil {
		return u.Log.Error(ctx, err)
	}

	return nil
}

func (u *UserRepository) Update(ctx context.Context) error {
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "UpdateUserRepository")
	defer span.End()

	switch ctx.Err() {
	case context.Canceled:
		return u.Log.Error(ctx, context.Canceled)
	case context.DeadlineExceeded:
		return u.Log.Error(ctx, context.DeadlineExceeded)
	default:
	}

	const q = `UPDATE users SET name = $1, updated_at = timezone('utc', now()), updated_by = $2 WHERE id = $3 RETURNING email`
	span.SetAttributes(attribute.String("db.query", q))
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		return u.Log.Error(ctx, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.UserEntity.Name, ctx.Value(myctx.Key("user_id")).(int64), u.UserEntity.ID).Scan(&u.UserEntity.Email)
	if err != nil {
		return u.Log.Error(ctx, err)
	}

	return nil
}

func (u *UserRepository) Delete(ctx context.Context) error {
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "DeleteUserRepository")
	defer span.End()

	switch ctx.Err() {
	case context.Canceled:
		return u.Log.Error(ctx, context.Canceled)
	case context.DeadlineExceeded:
		return u.Log.Error(ctx, context.DeadlineExceeded)
	default:
	}

	const q = `UPDATE users SET deleted_at = timezone('utc', now()), deleted_by = $1 WHERE id = $2`
	span.SetAttributes(attribute.String("db.query", q))
	span.SetAttributes(attribute.Int64("db.id", u.UserEntity.ID))

	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		return u.Log.Error(ctx, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ctx.Value(myctx.Key("user_id")).(int64), u.UserEntity.ID)
	if err != nil {
		return u.Log.Error(ctx, err)
	}

	return nil
}

func (u *UserRepository) List(ctx context.Context, search string) ([]model.User, error) {
	var list []model.User = make([]model.User, 0)
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "listUserRepository")
	defer span.End()

	switch ctx.Err() {
	case context.Canceled:
		return list, u.Log.Error(ctx, context.Canceled)
	case context.DeadlineExceeded:
		return list, u.Log.Error(ctx, context.DeadlineExceeded)
	default:
	}

	sb := strings.Builder{}
	sb.WriteString(`SELECT id, name, email FROM users WHERE deleted_at IS NULL`)
	var args []interface{}

	if len(search) > 0 {
		sb.WriteString(fmt.Sprintf(` AND name like %d`, len(args)+1))
		args = append(args, `%`+search+`%`)
	}
	span.SetAttributes(attribute.String("db.query", sb.String()))
	stmt, err := u.Db.PrepareContext(ctx, sb.String())
	if err != nil {
		return list, u.Log.Error(ctx, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return list, u.Log.Error(ctx, err)
	}

	defer rows.Close()

	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return list, u.Log.Error(ctx, err)
		}
		list = append(list, user)
	}

	if rows.Err() != nil {
		return list, u.Log.Error(ctx, rows.Err())
	}

	return list, nil
}

func (u *UserRepository) GetByEmail(ctx context.Context) error {
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "GetByEmailtUserRepository")
	defer span.End()

	switch ctx.Err() {
	case context.Canceled:
		return u.Log.Error(ctx, context.Canceled)
	case context.DeadlineExceeded:
		return u.Log.Error(ctx, context.DeadlineExceeded)
	default:
	}

	const q = `SELECT id, password FROM users WHERE email=$1 AND deleted_at IS NULL`
	span.SetAttributes(attribute.String("db.query", q))
	span.SetAttributes(attribute.String("db.email", u.UserEntity.Email))

	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		return u.Log.Error(ctx, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.UserEntity.Email).Scan(&u.UserEntity.ID, &u.UserEntity.Password)
	if err != nil {
		return u.Log.Error(ctx, err)
	}
	return nil
}
