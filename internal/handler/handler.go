package handler

import (
	"github.com/go-redis/redis/v7"
	"github.com/juicyluv/astral/internal/store"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Handler struct {
	router *httprouter.Router
	logger *zap.SugaredLogger
	redis  *redis.Client
	store  store.Store
}

type jsonResponse map[string]interface{}

func NewHandler(logger *zap.SugaredLogger, store store.Store, redis *redis.Client) *Handler {
	h := &Handler{
		router: httprouter.New(),
		logger: logger,
		store:  store,
		redis:  redis,
	}

	h.initRoutes()

	return h
}

func (h *Handler) GetRouter() *httprouter.Router {
	return h.router
}
