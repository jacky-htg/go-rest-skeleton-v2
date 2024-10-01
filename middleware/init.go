package middleware

import (
	"rest-skeleton/pkg/database"
	"rest-skeleton/pkg/logger"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/metric"
)

type Middleware struct {
	Log           *logger.Logger
	DB            *database.Database
	LatencyMetric metric.Int64Histogram
}

func (m *Middleware) WrapMiddleware(mw []func(httprouter.Handle) httprouter.Handle, handler httprouter.Handle) httprouter.Handle {

	for i := 0; i < len(mw); i++ {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
