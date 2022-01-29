package handler

import "net/http"

func (h *Handler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.IsTokenExpired(r)
		if err != nil {
			h.UnauthorizedResponse(w, r)
			return
		}

		next(w, r)
	}
}
