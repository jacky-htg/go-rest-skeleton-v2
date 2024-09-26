package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rest-skeleton/pkg/config"
	"rest-skeleton/pkg/database"
	"rest-skeleton/pkg/logger"
	"rest-skeleton/pkg/redis"
	"rest-skeleton/route"

	_ "github.com/lib/pq"
)

func main() {
	log := (&logger.Logger{}).New()

	if _, ok := os.LookupEnv("APP_NAME"); !ok {
		if err := config.Setup(".env"); err != nil {
			log.Error.Fatal(err)
		}
	}

	log.Info.Println("Starting Server at : ", os.Getenv("APP_PORT"))

	db, err := database.NewDatabase()
	if err != nil {
		log.Error.Fatalf("Could not connect to database: %s\n", err)
	}
	defer db.Conn.Close()

	redisClient, err := redis.NewCache(context.Background(), os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"), 24*time.Hour)
	if err != nil {
		log.Error.Fatalf("Could not connect to redis: %s\n", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error.Printf("Could not close redis connection: %s\n", err)
		}
	}()

	srv := &http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 30,
		Handler:      route.InitRoute(log, db, redisClient),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error.Printf("Server shutdown: %s\n", err)
	}

	log.Info.Println("Server exiting")
}
