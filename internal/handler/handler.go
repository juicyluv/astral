package handler

import "github.com/julienschmidt/httprouter"

type Handler struct {
	router *httprouter.Router
}

func NewHandler() *Handler {
	h := &Handler{
		router: httprouter.New(),
	}

	h.initRoutes()

	return h
}

func (h *Handler) GetRouter() *httprouter.Router {
	return h.router
}
