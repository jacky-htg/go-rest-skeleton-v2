package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"rest-skeleton/model"
	"rest-skeleton/pkg/logger"
	"rest-skeleton/pkg/myctx"
)

type UserRepository struct {
	Db         *sql.DB
	Log        *logger.Logger
	UserEntity model.User
}

func (u *UserRepository) Find(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		err := errors.New("request is canceled")
		u.Log.Error.Println(err)
		return err
	case context.DeadlineExceeded:
		err := errors.New("deadline is exceeded")
		u.Log.Error.Println(err)
		return err
	default:
	}

	const q = `SELECT id, name, email, password FROM users WHERE id=$1 AND deleted_at IS NULL`
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error.Println("error get users", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.UserEntity.ID).Scan(&u.UserEntity.ID, &u.UserEntity.Name, &u.UserEntity.Email, &u.UserEntity.Password)
	if err != nil {
		u.Log.Error.Println("error get users", err)
		return err
	}
	return nil
}

func (u *UserRepository) Save(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		err := errors.New("request is canceled")
		u.Log.Error.Println(err)
		return err
	case context.DeadlineExceeded:
		err := errors.New("deadline is exceeded")
		u.Log.Error.Println(err)
		return err
	default:
	}

	const q = `INSERT INTO users (name, password, email, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error.Println("error get users", err)
		return err
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
		u.Log.Error.Println("error create users", err)
		return err
	}

	return nil
}

func (u *UserRepository) Update(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		err := errors.New("request is canceled")
		u.Log.Error.Println(err)
		return err
	case context.DeadlineExceeded:
		err := errors.New("deadline is exceeded")
		u.Log.Error.Println(err)
		return err
	default:
	}

	const q = `UPDATE users SET name = $1, updated_at = timezone('utc', now()), updated_by = $2 WHERE id = $3 RETURNING email`
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error.Println("error get users", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.UserEntity.Name, ctx.Value(myctx.Key("user_id")).(int64), u.UserEntity.ID).Scan(&u.UserEntity.Email)
	if err != nil {
		u.Log.Error.Println("error update users", err)
		return err
	}

	return nil
}

func (u *UserRepository) Delete(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		err := errors.New("request is canceled")
		u.Log.Error.Println(err)
		return err
	case context.DeadlineExceeded:
		err := errors.New("deadline is exceeded")
		u.Log.Error.Println(err)
		return err
	default:
	}

	const q = `UPDATE users SET deleted_at = timezone('utc', now()), deleted_by = $1 WHERE id = $2`
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error.Println("error get users", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ctx.Value(myctx.Key("user_id")).(int64), u.UserEntity.ID)
	if err != nil {
		u.Log.Error.Println("error delete users", err)
		return err
	}

	return nil
}

func (u *UserRepository) List(ctx context.Context, search string) ([]model.User, error) {
	var list []model.User = make([]model.User, 0)

	switch ctx.Err() {
	case context.Canceled:
		err := errors.New("request is canceled")
		u.Log.Error.Println(err)
		return list, err
	case context.DeadlineExceeded:
		err := errors.New("deadline is exceeded")
		u.Log.Error.Println(err)
		return list, err
	default:
	}

	sb := strings.Builder{}
	sb.WriteString(`SELECT id, name, email FROM users WHERE deleted_at IS NULL`)
	var args []interface{}

	if len(search) > 0 {
		sb.WriteString(fmt.Sprintf(` AND name like %d`, len(args)+1))
		args = append(args, `%`+search+`%`)
	}
	stmt, err := u.Db.PrepareContext(ctx, sb.String())
	if err != nil {
		u.Log.Error.Println("error get users", err)
		return list, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		u.Log.Error.Println("error get users", err)
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			u.Log.Error.Println("error get users", err)
			return list, err
		}
		list = append(list, user)
	}

	if rows.Err() != nil {
		u.Log.Error.Println("error get users", err)
		return list, err
	}

	return list, nil
}

func (u *UserRepository) GetByEmail(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		err := errors.New("request is canceled")
		u.Log.Error.Println(err)
		return err
	case context.DeadlineExceeded:
		err := errors.New("deadline is exceeded")
		u.Log.Error.Println(err)
		return err
	default:
	}

	const q = `SELECT id, password FROM users WHERE email=$1 AND deleted_at IS NULL`
	stmt, err := u.Db.PrepareContext(ctx, q)
	if err != nil {
		u.Log.Error.Println("error get user by email", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.UserEntity.Email).Scan(&u.UserEntity.ID, &u.UserEntity.Password)
	if err != nil {
		u.Log.Error.Println("error get user by email", err)
		return err
	}
	return nil
}
