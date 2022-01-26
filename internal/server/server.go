package server

import (
	"net/http"

	"github.com/juicyluv/astral/internal/handler"
	"github.com/juicyluv/astral/internal/store"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	cfg    *Config
	logger *zap.SugaredLogger
	db     store.Store
}

func NewServer(cfg *Config, logger *zap.SugaredLogger, store store.Store) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		db:     store,
		server: &http.Server{
			Addr:           cfg.Port,
			WriteTimeout:   cfg.WriteTimeout,
			ReadTimeout:    cfg.ReadTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
			Handler:        handler.NewHandler().GetRouter(),
		},
	}
}

func (s *Server) Run() error {
	s.logger.Infof("Server started on port %s\n", s.cfg.Port)
	return s.server.ListenAndServe()
}
