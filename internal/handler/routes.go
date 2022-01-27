package handler

import "net/http"

func (h *Handler) initRoutes() {
	h.router.NotFound = http.HandlerFunc(h.notFoundResponse)
	h.router.MethodNotAllowed = http.HandlerFunc(h.methodNotAllowedResponse)

	h.router.HandlerFunc(http.MethodGet, "/api/health", h.health)

	h.router.HandlerFunc(http.MethodPost, "/api/user", h.createUser)
	h.router.HandlerFunc(http.MethodGet, "/api/user", h.listUser)
	h.router.HandlerFunc(http.MethodGet, "/api/user/:id", h.getUser)
	h.router.HandlerFunc(http.MethodPut, "/api/user/:id", h.updateUser)
	h.router.HandlerFunc(http.MethodDelete, "/api/user/:id", h.deleteUser)

	h.router.HandlerFunc(http.MethodGet, "/api/post", h.listPost)
	h.router.HandlerFunc(http.MethodPost, "/api/post", h.createPost)
	h.router.HandlerFunc(http.MethodGet, "/api/post/:id", h.getPost)
	h.router.HandlerFunc(http.MethodPut, "/api/post/:id", h.updatePost)
	h.router.HandlerFunc(http.MethodDelete, "/api/post/:id", h.deletePost)
}
