package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (m *Middleware) PanicRecovery(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				stack := strings.Builder{}
				for i := 1; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					stack.WriteString(fmt.Sprintf("  %s:%d\n", file, line))
				}

				m.Log.Error(r.Context(), fmt.Errorf("%v\n panic: %s", err, stack.String()))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r, ps)
	})
}
