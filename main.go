package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rest-skeleton/internal/pkg/config"
	"rest-skeleton/internal/pkg/database"
	"rest-skeleton/internal/pkg/logger"
	"rest-skeleton/internal/pkg/redis"
	"rest-skeleton/internal/pkg/telemetry"
	"rest-skeleton/internal/route"

	_ "github.com/lib/pq"
)

// @title Rest Skeleton API
// @version 1.0
// @description This is a sample server API.
// @Schemes http
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	if _, ok := os.LookupEnv("APP_NAME"); !ok {
		if err := config.Setup(".env"); err != nil {
			fmt.Printf("failed to setup config: %v", err)
			os.Exit(1)
		}
	}

	meter, err := telemetry.NewMeter(context.Background())
	if err != nil {
		fmt.Println("failed to create meter", err)
		os.Exit(1)
	}

	// go telemetry.CollectMachineResourceMetrics(meter) // dihapus karena redundant dengan matric go secara umum

	today := time.Now().Format("2006-01-02")
	logFileName := "log/api-" + today + ".log"
	log := logger.New(logFileName)

	fmt.Println("Starting Server at : "+os.Getenv("APP_PORT"), "")

	shutdown, err := telemetry.InitTracing()
	if err != nil {
		fmt.Printf("failed to initialize tracing: %v", err)
		os.Exit(1)
	}
	defer shutdown(context.Background())

	latencyMetric, errorCountMetric, err := telemetry.SetMetric(meter)
	if err != nil {
		fmt.Printf("failed to initialize metrics: %v", err)
		os.Exit(1)
	}
	log.ErrorCountMetric = errorCountMetric

	db, err := database.NewDatabase()
	if err != nil {
		fmt.Printf("Could not connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Conn.Close()

	redisClient, err := redis.NewCache(context.Background(), os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"), 24*time.Hour)
	if err != nil {
		fmt.Printf("Could not connect to redis: %v", err)
	}
	defer redisClient.Close()

	srv := &http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 30,
		Handler:      route.ApiRoute(log, db, redisClient, latencyMetric),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println("listen and serve", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server ...", "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server shutdown", err)
	}

	fmt.Println("Server exiting")
}
