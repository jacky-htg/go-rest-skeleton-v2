// ./middleware/semaphore.go
package middleware

import (
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// SemaphoreMiddleware is a middleware that limits the number of concurrent requests
func (m *Middleware) Semaphore(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		maxConcurrentRequests, err := strconv.Atoi(os.Getenv("CONCURRENCY_LIMIT"))
		if err != nil {
			m.Log.Error(r.Context(), err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		sem := make(chan struct{}, maxConcurrentRequests) // Create a semaphore with a maximum capacity

		// Acquire a semaphore slot
		sem <- struct{}{}

		// Ensure the semaphore is released after the request is processed
		defer func() {
			<-sem
		}()

		next(w, r, ps)
	}
}
