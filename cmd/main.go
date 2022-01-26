package main

import (
	"flag"
	"log"

	"github.com/juicyluv/astral/configs"
	"github.com/juicyluv/astral/internal/server"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config-path", "configs/dev.yml", "the application config path")
)

func main() {
	flag.Parse()

	prodLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer prodLogger.Sync()
	logger := prodLogger.Sugar()

	if err := configs.LoadConfigs(*configPath); err != nil {
		logger.Fatal(err)
	}

	config := server.NewConfig(*configPath)

	server := server.NewServer(&config, logger)

	if err := server.Run(); err != nil {
		logger.Fatal(err)
	}
}
