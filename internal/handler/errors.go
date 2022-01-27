package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4"
)

var (
	errNoRows         = pgx.ErrNoRows
	errNoRowsResponse = errors.New("record not found")
)

// errorResponse logs an error and sends a JSON response with a given status code.
func (h *Handler) errorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message interface{}) {
	msg := jsonResponse{"error": message}

	if err := sendJSON(w, msg, statusCode, nil); err != nil {
		h.logError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// internalErrorResponse logs the error message and sends a 500 Internal Server Error by using
// errorResponse helper function.
func (h *Handler) internalErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(err)

	message := "the server encountered a problem and could not process your request"
	h.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse sends 404 Not Found status code and
// JSON error response to the client.
func (h *Handler) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	h.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse sends a 405 Method Not Allowed
// status code and JSON response to the client.
func (h *Handler) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	h.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse sends a 404 Bad Request status code and JSON error message.
func (h *Handler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// recordNotFoundResponse sends a 404 Not Found response when record not found in storage.
func (h *Handler) recordNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	h.errorResponse(w, r, http.StatusBadRequest, errNoRowsResponse.Error())
}

func (h *Handler) invalidRequestBodyResponse(w http.ResponseWriter, r *http.Request) {
	h.errorResponse(w, r, http.StatusUnprocessableEntity, "invalid request body")
}
