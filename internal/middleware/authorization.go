// ./middleware/authentication.go
package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"rest-skeleton/internal/pkg/myctx"
	"rest-skeleton/internal/repository"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (m *Middleware) Authorization(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		path := r.URL.Path
		for _, param := range ps {
			path = strings.Replace(path, "/"+ps.ByName(param.Key), "/:"+param.Key, 1)
		}
		ctx := context.WithValue(r.Context(), myctx.Key("path"), path)
		r = r.WithContext(ctx)

		authRepository := repository.AuthRepository{Db: m.DB.Conn, Log: m.Log}
		hasAuth, err := authRepository.HasAuth(r.Context(), r.Method+" "+path)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !hasAuth {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r, ps)
	})
}
