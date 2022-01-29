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
	"github.com/juicyluv/astral/internal/model"
)

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var login model.Auth

	if err := readJSON(w, r, &login); err != nil {
		h.invalidRequestBodyResponse(w, r)
		return
	}

	login.Email = strings.ToLower(login.Email)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	user, err := h.store.User().FindByEmail(ctx, login.Email)
	if err != nil {
		if errors.Is(err, errNoRowsResponse) {
			h.notFoundResponse(w, r)
		} else {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	if !user.ComparePassword(login.Password) || user.Email != login.Email {
		h.badRequestResponse(w, r, errors.New("invalid email or password"))
		return
	}

	user.ClearPassword()

	tokens, err := login.CreateToken(user.Id)
	if err != nil {
		h.internalErrorResponse(w, r, err)
		return
	}

	err = h.saveTokenInformation(user.Id, tokens)
	if err != nil {
		h.internalErrorResponse(w, r, err)
		return
	}

	tokenResponse := jsonResponse{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	}

	if err := sendJSON(w, &tokenResponse, http.StatusOK, nil); err != nil {
		h.internalErrorResponse(w, r, err)
	}
}

func (h *Handler) saveTokenInformation(userId int, td *model.TokenDetails) error {
	// Converting Unix to UTC
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	// Set Access token
	err := h.redis.Set(td.AccessUuid, strconv.Itoa(userId), at.Sub(now)).Err()
	if err != nil {
		return err
	}

	// Set Refresh token
	err = h.redis.Set(td.RefreshUuid, strconv.Itoa(userId), rt.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) ExtractToken(r *http.Request) (string, error) {
	err := errors.New("invalid authorization header")

	bearer := r.Header.Get("Authorization")
	s := strings.Split(bearer, " ")
	if len(s) != 2 {
		return "", err
	}

	if s[0] != "Bearer" {
		return "", err
	}

	return s[1], nil
}

func (h *Handler) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString, err := h.ExtractToken(r)
	if err != nil {
		return nil, err
	}

	// Check token signing method to be HMAC
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected token signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// IsTokenExiped will check whether token valid or has already expired
func (h *Handler) IsTokenExpired(r *http.Request) error {
	token, err := h.VerifyToken(r)
	if err != nil {
		return err
	}

	if !token.Valid {
		return err
	}

	return nil
}

// GetTokenMetadata will extract token metadata and return it if there is no error
func (h *Handler) GetTokenMetadata(r *http.Request) (*model.TokenMetadata, error) {
	token, err := h.VerifyToken(r)
	if err != nil {
		return nil, err
	}

	err = errors.New("token is not valid")
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			fmt.Println("access")
			return nil, err
		}

		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			fmt.Println("user id int")
			return nil, err
		}

		return &model.TokenMetadata{
			AccessUuid: accessUuid,
			UserId:     int(userId),
		}, nil
	}

	return nil, err
}

// FetchTokenDataFromRedis will return the user id depends on given token metadata
func (h *Handler) FetchTokenDataFromRedis(details *model.TokenMetadata) (int, error) {
	userIdString, err := h.redis.Get(details.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
