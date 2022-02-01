package handler

import (
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/juicyluv/astral/internal/store"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Handler struct {
	router *httprouter.Router
	logger *zap.SugaredLogger
	redis  *redis.Client
	store  store.Store

	requestTimeout time.Duration
}

type jsonResponse map[string]interface{}

func NewHandler(logger *zap.SugaredLogger, store store.Store, redis *redis.Client) *Handler {
	h := &Handler{
		router: httprouter.New(),
		logger: logger,
		store:  store,
		redis:  redis,

		requestTimeout: time.Duration(viper.GetInt("http.requestTimeout")) * time.Second,
	}

	h.initRoutes()

	return h
}

// GetRouter returns a router instance pointer
func (h *Handler) GetRouter() *httprouter.Router {
	return h.router
}
