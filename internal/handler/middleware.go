package handler

import "net/http"

// RequireAuth middleware will check if token is presented and if it is valid
func (h *Handler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.isTokenExpired(r)
		if err != nil {
			h.unauthorizedResponse(w, r)
			return
		}

		next(w, r)
	}
}
