package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/juicyluv/astral/internal/model"
)

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	token, err := h.getTokenMetadata(r)
	if err != nil {
		h.unauthorizedResponse(w, r)
		return
	}

	userId, err := h.fetchTokenDataFromRedis(token)
	if err != nil {
		h.unauthorizedResponse(w, r)
		return
	}

	var post model.Post

	if err := readJSON(w, r, &post); err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	post.Author.Id = userId

	if err := post.Validate(); err != nil {
		h.errorResponse(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	postId, err := h.store.Post().Create(ctx, &post)
	if err != nil {
		h.internalErrorResponse(w, r, err)
		return
	}

	err = sendJSON(w, jsonResponse{"id": postId}, http.StatusOK, nil)
	if err != nil {
		h.internalErrorResponse(w, r, err)
	}
}

func (h *Handler) listPost(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	posts, err := h.store.Post().FindAll(ctx)
	if err != nil {
		if errors.Is(err, errNoRows) {
			h.recordNotFoundResponse(w, r)
		} else {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	err = sendJSON(w, posts, http.StatusOK, nil)
	if err != nil {
		h.internalErrorResponse(w, r, err)
	}
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	postId, err := readIdParam(r)
	if err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.store.Post().FindById(ctx, int(postId))
	if err != nil {
		if errors.Is(err, errNoRows) {
			h.recordNotFoundResponse(w, r)
		} else {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	err = sendJSON(w, post, http.StatusOK, nil)
	if err != nil {
		h.internalErrorResponse(w, r, err)
	}
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	postId, err := readIdParam(r)
	if err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var post model.UpdatePostDto

	if err := readJSON(w, r, &post); err != nil {
		h.invalidRequestBodyResponse(w, r)
		return
	}

	if err := post.Validate(); err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// If updating post's author, check if author with this id exists
	if post.AuthorId != nil {
		_, err = h.store.User().FindById(ctx, *post.AuthorId)
		if err != nil {
			if errors.Is(err, errNoRows) {
				h.errorResponse(w, r, http.StatusBadRequest, "there is no user with this id")
			} else {
				h.internalErrorResponse(w, r, err)
			}
			return
		}
	}

	err = h.store.Post().Update(ctx, int(postId), &post)
	if err != nil {
		if errors.Is(err, errNoRows) {
			h.recordNotFoundResponse(w, r)
		} else {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	err = sendJSON(w, nil, http.StatusOK, nil)
	if err != nil {
		h.internalErrorResponse(w, r, err)
	}
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	postId, err := readIdParam(r)
	if err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.store.Post().Delete(ctx, int(postId))
	if err != nil {
		if errors.Is(err, errNoRows) {
			h.recordNotFoundResponse(w, r)
		} else {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	err = sendJSON(w, nil, http.StatusOK, nil)
	if err != nil {
		h.internalErrorResponse(w, r, err)
	}
}
