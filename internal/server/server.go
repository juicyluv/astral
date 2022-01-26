package server

import (
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	cfg    *Config
	logger *zap.SugaredLogger
}

func NewServer(cfg *Config, logger *zap.SugaredLogger) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		server: &http.Server{
			Addr:           cfg.Port,
			WriteTimeout:   cfg.WriteTimeout,
			ReadTimeout:    cfg.ReadTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
		},
	}
}

func (s *Server) Run() error {
	s.logger.Infof("Server started on port %s\n", s.cfg.Port)
	return s.server.ListenAndServe()
}
