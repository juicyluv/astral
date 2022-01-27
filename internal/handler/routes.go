package handler

import "net/http"

func (h *Handler) initRoutes() {
	h.router.HandlerFunc(http.MethodGet, "/api/health", h.health)
	h.router.HandlerFunc(http.MethodPost, "/api/user", h.createUser)
	h.router.HandlerFunc(http.MethodGet, "/api/user", h.listUser)
	h.router.HandlerFunc(http.MethodGet, "/api/user/:id", h.getUser)
	h.router.HandlerFunc(http.MethodDelete, "/api/user/:id", h.deleteUser)
}
