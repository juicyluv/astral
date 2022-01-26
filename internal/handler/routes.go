package handler

import "net/http"

func (h *Handler) initRoutes() {
	h.router.HandlerFunc(http.MethodGet, "/api/health", h.health)
}
