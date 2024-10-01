package route

import (
	"fmt"
	"net/http"
	"os"
	_ "rest-skeleton/docs"
	"rest-skeleton/internal/handler"
	"rest-skeleton/internal/middleware"
	"rest-skeleton/internal/pkg/database"
	"rest-skeleton/internal/pkg/logger"
	"rest-skeleton/internal/pkg/redis"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel/metric"
)

func ApiRoute(log *logger.Logger, db *database.Database, cache *redis.Cache, latencyMetric metric.Int64Histogram) *httprouter.Router {
	router := httprouter.New()
	router.ServeFiles("/docs/*filepath", http.Dir("./docs"))

	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s:%s/docs/swagger.json", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))),
	)
	router.Handler("GET", "/swagger/*filepath", swaggerHandler)
	router.Handler("GET", "/metrics", promhttp.Handler())

	var mid middleware.Middleware = middleware.Middleware{Log: log, DB: db, LatencyMetric: latencyMetric}
	publicMiddlewares := []func(httprouter.Handle) httprouter.Handle{
		mid.TraceAndMetricLatency,
		mid.CORS,
		mid.PanicRecovery,
		mid.Semaphore,
		mid.RateLimit,
	}
	privateMiddlewares := append(publicMiddlewares, mid.Authentication, mid.Authorization)

	userHandler := handler.Users{Log: log, DB: db.Conn, Cache: cache}
	authHandler := handler.Auths{Log: log, DB: db.Conn}

	router.POST("/login", mid.WrapMiddleware(publicMiddlewares, authHandler.Login))
	router.GET("/users", mid.WrapMiddleware(privateMiddlewares, userHandler.List))
	router.GET("/users/:id", mid.WrapMiddleware(privateMiddlewares, userHandler.GetById))
	router.POST("/users", mid.WrapMiddleware(privateMiddlewares, userHandler.Create))
	router.PUT("/users/:id", mid.WrapMiddleware(privateMiddlewares, userHandler.Update))
	router.DELETE("/users/:id", mid.WrapMiddleware(privateMiddlewares, userHandler.Delete))

	return router
}
