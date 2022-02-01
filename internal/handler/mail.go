package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/juicyluv/astral/internal/mail"
	"github.com/spf13/viper"
)

func (h *Handler) confirmEmail(w http.ResponseWriter, r *http.Request) {
	errInvalidToken := errors.New("invalid token")

	// Parse token from URL query
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		h.badRequestResponse(w, r, errors.New("empty token"))
		return
	}

	// Parse and check token, get metadata
	// Check token signing method to be HMAC
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected token signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		h.badRequestResponse(w, r, errInvalidToken)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		h.badRequestResponse(w, r, errInvalidToken)
		return
	}

	userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		h.badRequestResponse(w, r, errInvalidToken)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	err = h.store.User().ConfirmEmail(ctx, int(userId))
	if err != nil {
		h.internalErrorResponse(w, r, err)
		return
	}

	user, err := h.store.User().FindById(ctx, int(userId))
	if err != nil {
		h.internalErrorResponse(w, r, err)
		return
	}

	message := "You successfully verified your account."
	err = mail.SendEmail(strings.ToLower(user.Email), message)
	if err != nil {
		h.internalErrorResponse(w, r, errors.New("could not send email"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) createEmailToken(userId int) (string, error) {
	tokenExpTimeDays := time.Duration(viper.GetInt("mail.tokenExpTime")) * time.Hour * 24

	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["exp"] = tokenExpTimeDays

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}
