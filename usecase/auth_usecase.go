package usecase

import (
	"context"
	"errors"
	"net/http"
	"rest-skeleton/dto"
	"rest-skeleton/model"
	"rest-skeleton/pkg/database"
	"rest-skeleton/pkg/jwttoken"
	"rest-skeleton/pkg/logger"
	"rest-skeleton/repository"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthUC struct {
	Log *logger.Logger
	DB  *database.Database
}

func (uc AuthUC) Login(ctx context.Context, loginRequest dto.LoginRequest) (string, int, error) {
	switch ctx.Err() {
	case context.Canceled:
		err := errors.New("request is canceled")
		uc.Log.Error.Println(err)
		return "", http.StatusInternalServerError, err
	case context.DeadlineExceeded:
		err := errors.New("deadline is exceeded")
		uc.Log.Error.Println(err)
		return "", http.StatusInternalServerError, err
	default:
	}

	userRepo := repository.UserRepository{Log: uc.Log, Db: uc.DB.Conn, UserEntity: model.User{Email: loginRequest.Email}}
	if err := userRepo.GetByEmail(ctx); err != nil {
		return "", http.StatusInternalServerError, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(strings.TrimSpace(userRepo.UserEntity.Password)), []byte(loginRequest.Password)); err != nil {
		uc.Log.Error.Println("invalid credentials", err)
		return "", http.StatusUnauthorized, err
	}

	token, err := jwttoken.ClaimToken(userRepo.UserEntity.Email)
	if err != nil {
		uc.Log.Error.Println(err)
		return "", http.StatusInternalServerError, err
	}

	return token, http.StatusOK, nil
}
