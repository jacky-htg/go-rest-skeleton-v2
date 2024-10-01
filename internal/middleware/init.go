package middleware

import (
	"database/sql"
	"rest-skeleton/internal/pkg/logger"
	"rest-skeleton/internal/pkg/redis"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/metric"
)

type Middleware struct {
	Log           *logger.Logger
	DB            *sql.DB
	Cache         *redis.Cache
	LatencyMetric metric.Int64Histogram
}

func (m *Middleware) WrapMiddleware(mw []func(httprouter.Handle) httprouter.Handle, handler httprouter.Handle) httprouter.Handle {

	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
