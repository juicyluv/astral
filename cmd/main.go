package main

import (
	"context"
	"flag"
	"log"

	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/astral/configs"
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

	store := postgres.NewPostgres(conn, logger)

	server := server.NewServer(&config, logger, store, redis)

	// Run the server
	if err := server.Run(); err != nil {
		logger.Fatal(err)
	}
}
