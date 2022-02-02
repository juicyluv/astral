package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/astral/configs"
	"github.com/juicyluv/astral/internal/queue"
	"github.com/juicyluv/astral/internal/server"
	"github.com/juicyluv/astral/internal/store/postgres"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config-path", "configs/dev.yml", "the application config path")
)

func main() {
	flag.Parse()

	// Logger initialization
	prodLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer prodLogger.Sync()
	logger := prodLogger.Sugar()

	// Load config files
	if err := configs.LoadConfigs(*configPath); err != nil {
		logger.Fatal(err)
	}

	// Create config instance
	config := server.NewConfig(*configPath)

	// Create database connection
	conn, err := pgx.Connect(context.Background(), config.DbDSN)
	if err != nil {
		logger.Fatal(err)
	}

	// Try to connect to database
	if err = conn.Ping(context.Background()); err != nil {
		logger.Fatal(err)
	}
	logger.Info("connected to database")

	redis := redis.NewClient(&redis.Options{
		Addr: config.RedisDSN,
	})

	if _, err = redis.Ping().Result(); err != nil {
		logger.Fatal(err)
	}
	logger.Info("cache has been connected")

	queue, err := queue.NewQueue(logger, queue.NewConfig())
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("queue has been connected")

	store := postgres.NewPostgres(conn, logger)

	server := server.NewServer(&config, logger, store, redis, queue)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Run the server
	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	<-ctx.Done()

	logger.Info("shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := store.Close(context.Background()); err != nil {
		logger.Errorf("error occured on db connection close: %s", err.Error())
	}
	logger.Info("database has been closed")

	if err := redis.Close(); err != nil {
		logger.Errorf("error occured on redis connection close: %s", err.Error())
	}
	logger.Info("redis has been closed")

	if err := queue.Close(); err != nil {
		logger.Errorf("error occured on queue connection close: %s", err.Error())
	}
	logger.Info("queue has been closed")

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("shutdown: %w", err)
	}

	longShutdown := make(chan struct{}, 1)

	go func() {
		time.Sleep(3 * time.Second)
		longShutdown <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		logger.Errorf("server shutdown: %w", ctx.Err())
	case <-longShutdown:
		logger.Info("server has been shutted down")
	}
}
