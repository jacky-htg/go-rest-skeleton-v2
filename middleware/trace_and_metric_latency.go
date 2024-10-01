package middleware

import (
	"context"
	"net/http"
	"os"
	"rest-skeleton/pkg/myctx"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (m *Middleware) TraceAndMetricLatency(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()

		ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(r.Context(), r.URL.Path)
		defer span.End()

		traceID := span.SpanContext().TraceID().String()
		ctx = context.WithValue(ctx, myctx.Key("traceID"), traceID)

		rw := &responseWriter{w, http.StatusOK}
		next(rw, r.WithContext(ctx), ps)

		duration := time.Since(start).Milliseconds()
		myAttribute := []attribute.KeyValue{
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.Path),
			attribute.Int("http.status_code", rw.statusCode),
		}

		m.LatencyMetric.Record(ctx, duration, metric.WithAttributeSet(attribute.NewSet(myAttribute...)))
		myAttribute = append(myAttribute, attribute.Float64("http.duration_ms", float64(duration)))
		span.SetAttributes(myAttribute...)
	})
}
