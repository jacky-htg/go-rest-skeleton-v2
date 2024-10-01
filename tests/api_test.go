package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"rest-skeleton/internal/pkg/database"
	"rest-skeleton/internal/pkg/migration"
	"rest-skeleton/internal/pkg/redis"
	"sync"
	"time"

	redisDriver "github.com/go-redis/redis/v8"
)

var (
	once          sync.Once
	dbInstance    *sql.DB
	unitTeardown  func()
	dbTeardown    func()
	redisTeardown func()
	redisInstance *redis.Cache
)

func NewUnit() (*sql.DB, *redis.Cache, func()) {
	// Start container and initialize the database only once
	once.Do(func() {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			redisInstance, redisTeardown = NewRedis()
		}()

		go func() {
			defer wg.Done()
			dbInstance, dbTeardown = NewPosgrest()

			if err := migration.Migrate(dbInstance); err != nil {
				dbInstance.Close()
				dbTeardown()
				fmt.Println("migrating", err)
				os.Exit(1)
			} else {
				fmt.Println("Migrated database successfully")
			}
		}()

		wg.Wait()

		unitTeardown = func() {
			dbInstance.Close()
			redisInstance.Close()
			dbTeardown()
			redisTeardown()
		}
	})

	return dbInstance, redisInstance, unitTeardown
}

func NewPosgrest() (*sql.DB, func()) {
	database.StartPostgresContainer()
	psqlInfo := "host=localhost port=54320 user=postgres password=1234 dbname=postgres sslmode=disable"

	var err error
	dbInstance, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		database.StopPostgresContainer()
		fmt.Println("opening database connection", err)
		return nil, func() {}
	}

	fmt.Println("waiting for database to be ready")

	// Wait for the database to be ready.
	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = dbInstance.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 1000 * time.Millisecond)
	}

	if pingError != nil {
		database.StopPostgresContainer()
		fmt.Println(pingError)
	}

	return dbInstance, dbTeardown
}

func NewRedis() (*redis.Cache, func()) {
	database.StartRedisContainer()

	client := redisDriver.NewClient(&redisDriver.Options{
		Addr:     "localhost:63790",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	fmt.Println("waiting for redis to be ready")

	// Wait for the database to be ready.
	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = client.Ping(context.Background()).Err()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 1000 * time.Millisecond)
	}

	if pingError != nil {
		database.StopRedisContainer()
		fmt.Println(pingError)
	}

	redisTeardown := func() {
		database.StopRedisContainer()
	}

	return redis.NewCacheWithClient(context.Background(), client, 24*time.Hour), redisTeardown
}
