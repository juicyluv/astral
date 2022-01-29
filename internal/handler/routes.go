package handler

import "net/http"

func (h *Handler) initRoutes() {
	h.router.NotFound = http.HandlerFunc(h.notFoundResponse)
	h.router.MethodNotAllowed = http.HandlerFunc(h.methodNotAllowedResponse)

	h.router.HandlerFunc(http.MethodGet, "/api/health", h.health)

	h.router.HandlerFunc(http.MethodPost, "/api/auth/signin", h.login)
	h.router.HandlerFunc(http.MethodPost, "/api/auth/signup", h.createUser)

	h.router.HandlerFunc(http.MethodGet, "/api/users", h.listUser)
	h.router.HandlerFunc(http.MethodGet, "/api/users/:id", h.getUser)
	h.router.HandlerFunc(http.MethodPut, "/api/users/:id", h.updateUser)
	h.router.HandlerFunc(http.MethodDelete, "/api/users/:id", h.deleteUser)
	h.router.HandlerFunc(http.MethodGet, "/api/users/:id/posts", h.listUserPosts)

	h.router.HandlerFunc(http.MethodGet, "/api/posts", h.listPost)
	h.router.HandlerFunc(http.MethodPost, "/api/posts", h.createPost)
	h.router.HandlerFunc(http.MethodGet, "/api/posts/:id", h.getPost)
	h.router.HandlerFunc(http.MethodPut, "/api/posts/:id", h.updatePost)
	h.router.HandlerFunc(http.MethodDelete, "/api/posts/:id", h.deletePost)
}
