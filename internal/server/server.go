package server

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v7"
	"github.com/juicyluv/astral/internal/handler"
	"github.com/juicyluv/astral/internal/queue"
	"github.com/juicyluv/astral/internal/store"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	cfg    *Config
	logger *zap.SugaredLogger
	db     store.Store
}

func NewServer(cfg *Config, logger *zap.SugaredLogger, store store.Store, redis *redis.Client, queue *queue.Queue) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		db:     store,
		server: &http.Server{
			Addr:           cfg.Port,
			WriteTimeout:   cfg.WriteTimeout,
			ReadTimeout:    cfg.ReadTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
			Handler:        handler.NewHandler(logger, store, redis, queue).GetRouter(),
		},
	}
}

func (s *Server) Run() error {
	s.logger.Infof("Server started on port %s", s.cfg.Port)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
