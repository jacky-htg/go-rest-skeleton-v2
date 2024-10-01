package usecase

import (
	"context"
	"database/sql"
	"net/http"
	"rest-skeleton/dto"
	"rest-skeleton/model"
	"rest-skeleton/pkg/jwttoken"
	"rest-skeleton/pkg/logger"
	"rest-skeleton/repository"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthUC struct {
	Log *logger.Logger
	DB  *sql.DB
}

func (uc AuthUC) Login(ctx context.Context, loginRequest dto.LoginRequest) (string, int, error) {
	switch ctx.Err() {
	case context.Canceled:
		return "", http.StatusInternalServerError, uc.Log.Error(ctx, context.Canceled)
	case context.DeadlineExceeded:
		return "", http.StatusInternalServerError, uc.Log.Error(ctx, context.DeadlineExceeded)
	default:
	}

	userRepo := repository.UserRepository{Log: uc.Log, Db: uc.DB, UserEntity: model.User{Email: loginRequest.Email}}
	if err := userRepo.GetByEmail(ctx); err != nil {
		return "", http.StatusInternalServerError, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(strings.TrimSpace(userRepo.UserEntity.Password)), []byte(loginRequest.Password)); err != nil {
		return "", http.StatusUnauthorized, uc.Log.Error(ctx, err)
	}

	token, err := jwttoken.ClaimToken(userRepo.UserEntity.Email)
	if err != nil {
		return "", http.StatusInternalServerError, uc.Log.Error(ctx, err)
	}

	return token, http.StatusOK, nil
}
