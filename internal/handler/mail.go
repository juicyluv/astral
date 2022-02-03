package handler

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/juicyluv/astral/internal/mail"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// confirmEmail will parse URL to get confirmation token,
// then parses token to verify if it expired and get user id,
// and then updates the user as verified
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

	// Send email that user has been validated
	go func(logger *zap.SugaredLogger, email string) {
		subject := viper.GetString("mail.subject")
		buf, err := ioutil.ReadFile("./internal/mail/templates/confirm_success.html")
		if err != nil {
			h.logger.Error("cannot read html template")
			return
		}

		err = mail.SendEmail(email, subject, mail.MimeHTML, string(buf))
		if err != nil {
			logger.Infof("email was not sent: %v", err)
		}

		logger.Infof("email was sent to %s", email)
	}(h.logger, strings.ToLower(user.Email))

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
