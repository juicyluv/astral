package handler

import (
	"github.com/juicyluv/astral/internal/store"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Handler struct {
	router *httprouter.Router
	logger *zap.SugaredLogger
	store  store.Store
}

type jsonResponse map[string]interface{}

func NewHandler(logger *zap.SugaredLogger, store store.Store) *Handler {
	h := &Handler{
		router: httprouter.New(),
		logger: logger,
		store:  store,
	}

	h.initRoutes()

	return h
}

func (h *Handler) GetRouter() *httprouter.Router {
	return h.router
}
