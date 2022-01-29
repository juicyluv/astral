package handler

import "net/http"

func (h *Handler) initRoutes() {
	h.router.NotFound = http.HandlerFunc(h.notFoundResponse)
	h.router.MethodNotAllowed = http.HandlerFunc(h.methodNotAllowedResponse)

	// Helpers
	h.router.HandlerFunc(http.MethodGet, "/api/health", h.health)

	// Auth
	h.router.HandlerFunc(http.MethodPost, "/api/auth/signin", h.login)
	h.router.HandlerFunc(http.MethodPost, "/api/auth/signup", h.createUser)
	h.router.HandlerFunc(http.MethodGet, "/api/auth/signout", h.RequireAuth(h.logout))

	// Users
	h.router.HandlerFunc(http.MethodGet, "/api/users", h.listUser)
	h.router.HandlerFunc(http.MethodGet, "/api/users/:id", h.getUser)
	h.router.HandlerFunc(http.MethodPut, "/api/users/:id", h.RequireAuth(h.updateUser))
	h.router.HandlerFunc(http.MethodDelete, "/api/users/:id", h.RequireAuth(h.deleteUser))
	h.router.HandlerFunc(http.MethodGet, "/api/users/:id/posts", h.listUserPosts)

	// Posts
	h.router.HandlerFunc(http.MethodGet, "/api/posts", h.listPost)
	h.router.HandlerFunc(http.MethodPost, "/api/posts", h.RequireAuth(h.createPost))
	h.router.HandlerFunc(http.MethodGet, "/api/posts/:id", h.getPost)
	h.router.HandlerFunc(http.MethodPut, "/api/posts/:id", h.RequireAuth(h.updatePost))
	h.router.HandlerFunc(http.MethodDelete, "/api/posts/:id", h.RequireAuth(h.deletePost))
}
