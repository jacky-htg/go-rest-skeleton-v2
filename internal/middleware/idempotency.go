package middleware

import (
	"bytes"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (m *Middleware) Idempotency(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		idempotencyKey := r.Header.Get("Idempotency-Key")
		if (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete) && idempotencyKey == "" {
			http.Error(w, "Missing Idempotency-Key header", http.StatusBadRequest)
			return
		}

		if idempotencyKey == "" {
			next(w, r, ps)
			return
		}

		// Cek apakah kunci sudah ada di Redis
		ctx := r.Context()
		if cacheValue, isExist := m.Cache.Get(ctx, idempotencyKey); isExist {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(cacheValue.(string)))
			return
		}

		rw := &responseRecorder{ResponseWriter: w, body: new(bytes.Buffer)}
		next(rw, r, ps)

		m.Cache.SetTTL(10 * time.Minute)
		m.Cache.Add(ctx, idempotencyKey, rw.body.String())
		m.Cache.ResetTTL()
	})
}

// responseRecorder adalah struct untuk merekam respons
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseRecorder) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseRecorder) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}
