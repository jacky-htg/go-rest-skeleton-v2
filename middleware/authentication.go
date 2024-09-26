// ./middleware/authentication.go
package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"rest-skeleton/model"
	"rest-skeleton/pkg/jwttoken"
	"rest-skeleton/pkg/myctx"
	"rest-skeleton/repository"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (m *Middleware) Authentication(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, email := jwttoken.ValidateToken(token)
		if !isValid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userRepo := repository.UserRepository{Log: m.Log, Db: m.DB.Conn, UserEntity: model.User{Email: email}}
		if err := userRepo.GetByEmail(r.Context()); err != nil && err != sql.ErrNoRows {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else if err == sql.ErrNoRows {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), myctx.Key("email"), email)
		ctx = context.WithValue(ctx, myctx.Key("user_id"), userRepo.UserEntity.ID)
		r = r.WithContext(ctx)

		next(w, r, ps)
	})
}
