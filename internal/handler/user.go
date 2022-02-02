package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"text/template"

	"github.com/juicyluv/astral/internal/mail"
	"github.com/juicyluv/astral/internal/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := readJSON(w, r, &user)
	if err != nil {
		h.errorResponse(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = user.Validate(); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	found, err := h.store.User().FindByEmail(ctx, user.Email)
	if err != nil {
		if err != errNoRows {
			if found != nil {
				h.badRequestResponse(w, r, errors.New("email already taken"))
			} else {
				h.internalErrorResponse(w, r, err)
			}
			return
		}
	}

	if err = user.HashPassword(); err != nil {
		h.internalErrorResponse(w, r, err)
		return
	}

	userId, err := h.store.User().Create(ctx, &user)
	if err != nil {
		h.internalErrorResponse(w, r, err)
		return
	}

	token, err := h.createEmailToken(userId)
	if err != nil {
		h.internalErrorResponse(w, r, errors.New("could not create email token"))
		return
	}

	go func(logger *zap.SugaredLogger, username, email, token string) {
		subject := viper.GetString("mail.subject")
		filepath := "./internal/mail/templates/confirm_request.html"
		t := template.Must(template.ParseFiles(filepath))

		b := make([]byte, 0)
		buf := bytes.NewBuffer(b)

		t.ExecuteTemplate(buf, "confirm_request.html", struct {
			Username    string
			ConfirmLink string
		}{
			Username:    username,
			ConfirmLink: "http://localhost:8080/api/confirmation?token=" + token,
		})

		err = mail.SendEmail(email, subject, mail.MimeHTML, buf.String())
		if err != nil {
			logger.Infof("email was not sent: %v", err)
		}

		logger.Infof("email was sent to %s", email)
	}(h.logger, user.Username, user.Email, token)

	err = sendJSON(w, jsonResponse{"id": userId}, http.StatusOK, nil)
	if err != nil {
		h.internalErrorResponse(w, r, err)
	}
}

func (h *Handler) listUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
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

	err = sendJSON(w, users, http.StatusOK, nil)
	if err != nil {
		h.internalErrorResponse(w, r, err)
	}
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
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

	err = sendJSON(w, user, http.StatusOK, nil)
	if err != nil {
		h.internalErrorResponse(w, r, err)
	}
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	userId, err := readIdParam(r)
	if err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var user model.UpdateUserDto

	if err := readJSON(w, r, &user); err != nil {
		h.invalidRequestBodyResponse(w, r)
		return
	}

	if err := user.Validate(); err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.store.User().Update(ctx, int(userId), &user)
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

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
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

	err = sendJSON(w, nil, http.StatusOK, nil)
	if err != nil {
		h.internalErrorResponse(w, r, err)
	}
}

func (h *Handler) listUserPosts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	userId, err := readIdParam(r)
	if err != nil {
		h.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.store.User().FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, errNoRows) {
			h.errorResponse(w, r, http.StatusBadRequest, "user with this id not found")
		} else {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	posts, err := h.store.Post().FindUserPosts(ctx, userId)
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
