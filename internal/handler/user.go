package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/juicyluv/astral/internal/model"
)

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := readJSON(w, r, &user)
	if err != nil {
		h.errorResponse(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	userId, err := h.store.User().Create(ctx, &user)
	if err != nil {
		h.internalErrorResponse(w, r, err)
		return
	}

	sendJSON(w, jsonResponse{"id": userId}, http.StatusOK, nil)
}

func (h *Handler) listUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	users, err := h.store.User().FindAll(ctx)
	if err != nil {
		if errors.Is(err, errNoRows) {
			h.recordNotFoundResponse(w, r)
		} else {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	sendJSON(w, users, http.StatusOK, nil)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	userId, err := readIdParam(r)
	if err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.store.User().FindById(ctx, int(userId))
	if err != nil {
		if errors.Is(err, errNoRows) {
			h.recordNotFoundResponse(w, r)
		} else {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	sendJSON(w, user, http.StatusOK, nil)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	userId, err := readIdParam(r)
	if err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	err = h.store.User().Delete(ctx, int(userId))
	if err != nil {
		if errors.Is(err, errNoRows) {
			h.recordNotFoundResponse(w, r)
		} else {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	sendJSON(w, nil, http.StatusOK, nil)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	userId, err := readIdParam(r)
	if err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.store.User().Delete(ctx, int(userId))
	if err != nil {
		if errors.Is(err, errNoRows) {
			h.recordNotFoundResponse(w, r)
		} else {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	sendJSON(w, nil, http.StatusOK, nil)
}
