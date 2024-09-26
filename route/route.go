package route

import (
	"rest-skeleton/handler"
	"rest-skeleton/middleware"
	"rest-skeleton/pkg/database"
	"rest-skeleton/pkg/logger"
	"rest-skeleton/pkg/redis"

	"github.com/julienschmidt/httprouter"
)

func InitRoute(log *logger.Logger, db *database.Database, cache *redis.Cache) *httprouter.Router {
	router := httprouter.New()

	var mid middleware.Middleware = middleware.Middleware{Log: log, DB: db}
	publicMiddlewares := []func(httprouter.Handle) httprouter.Handle{
		mid.PanicRecovery,
		mid.CORS,
		mid.Semaphore,
		mid.RateLimit,
	}
	privateMiddlewares := append(publicMiddlewares, mid.Authentication, mid.Authorization)

	userHandler := handler.Users{Log: log, DB: db, Cache: cache}
	authHandler := handler.Auths{Log: log, DB: db}

	router.POST("/login", mid.WrapMiddleware(publicMiddlewares, authHandler.Login))
	router.GET("/users", mid.WrapMiddleware(privateMiddlewares, userHandler.List))
	router.GET("/users/:id", mid.WrapMiddleware(privateMiddlewares, userHandler.GetById))
	router.POST("/users", mid.WrapMiddleware(privateMiddlewares, userHandler.Create))
	router.PUT("/users/:id", mid.WrapMiddleware(privateMiddlewares, userHandler.Update))
	router.DELETE("/users/:id", mid.WrapMiddleware(privateMiddlewares, userHandler.Delete))

	return router
}
