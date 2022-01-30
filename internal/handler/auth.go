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
	"github.com/gofrs/uuid"
	"github.com/juicyluv/astral/internal/model"
	"github.com/spf13/viper"
)

// login parses user input, validates it and retrieves the user information
// from database. Then it creates a new pair of tokens and returns it
// to the client
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var login model.Auth

	if err := readJSON(w, r, &login); err != nil {
		h.invalidRequestBodyResponse(w, r)
		return
	}

	login.Email = strings.ToLower(login.Email)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
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

	tokens, err := h.createToken(user.Id)
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

// logout removes user session from cache and makes user tokens invalid
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	token, err := h.getTokenMetadata(r)
	if err != nil {
		h.unauthorizedResponse(w, r)
		return
	}

	deleted, err := h.removeUserTokenFromCache(token.AccessUuid)
	if err != nil || deleted == 0 {
		h.internalErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// createToken creates the jwt token and returns an error if something
// went wrong
func (h *Handler) createToken(userId int) (*model.TokenDetails, error) {
	// Access Token secret key
	accessSecret := os.Getenv("JWT_SECRET")
	if accessSecret == "" {
		accessSecret = "supersecretkey"
	}

	// Refresh token secret key
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = "supersecretkey"
	}

	// Access and Refresh token exp time
	tokenExpTimeMinutes := time.Duration(viper.GetInt("auth.tokenExpTime"))
	refreshExpTimeDays := time.Duration(viper.GetInt("auth.refreshExpTime"))

	// Token details structure
	td := model.TokenDetails{}

	// Generate access token uuid and exp time
	td.AtExpires = time.Now().Add(time.Minute * tokenExpTimeMinutes).Unix()
	accessUuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	td.AccessUuid = accessUuid.String()

	// Generate refresh token uuid and exp time
	td.RtExpires = time.Now().Add(time.Hour * 24 * refreshExpTimeDays).Unix()
	refreshUuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	td.RefreshUuid = refreshUuid.String()

	// Access token payload
	accessClaims := jwt.MapClaims{}
	accessClaims["user_id"] = userId
	accessClaims["access_uuid"] = td.AccessUuid
	accessClaims["authorized"] = true
	accessClaims["exp"] = td.AtExpires

	// Generate Access token
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := at.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}
	td.AccessToken = accessToken

	// Refresh token payload
	refreshClaims := jwt.MapClaims{}
	refreshClaims["user_id"] = userId
	refreshClaims["refresh_uuid"] = td.AccessUuid
	refreshClaims["authorized"] = true
	refreshClaims["exp"] = td.RtExpires

	// Generate Refresh token
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := rt.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}
	td.RefreshToken = refreshToken

	return &td, nil
}

// saveTokenInformation saves token information in the cache
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

// extractToken extracts token from Authorization request header
func (h *Handler) extractToken(r *http.Request) (string, error) {
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

// verifyToken checks whether token has HMAC encrypt algorithm
// and parses the given token
func (h *Handler) verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString, err := h.extractToken(r)
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
func (h *Handler) isTokenExpired(r *http.Request) error {
	token, err := h.verifyToken(r)
	if err != nil {
		return err
	}

	if !token.Valid {
		return err
	}

	return nil
}

// getTokenMetadata will extract token metadata and return it if there is no error
func (h *Handler) getTokenMetadata(r *http.Request) (*model.TokenMetadata, error) {
	token, err := h.verifyToken(r)
	if err != nil {
		return nil, err
	}

	err = errors.New("token is not valid")
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &model.TokenMetadata{
			AccessUuid: accessUuid,
			UserId:     int(userId),
		}, nil
	}

	return nil, err
}

// fetchTokenDataFromRedis will return the user id depends on given token metadata
func (h *Handler) fetchTokenDataFromRedis(details *model.TokenMetadata) (int, error) {
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

// removeUserTokenFromCache deleted the user information from cache
func (h *Handler) removeUserTokenFromCache(uuid string) (int, error) {
	deleted, err := h.redis.Del(uuid).Result()
	if err != nil {
		return 0, err
	}
	return int(deleted), nil
}

// refreshToken handles retrieving a new pair of tokens
func (h *Handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	type input struct {
		RefreshToken string `json:"refreshToken"`
	}

	tokenInput := input{}

	if err := readJSON(w, r, &tokenInput); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Parse token and make sure that it uses HMAC hash algorithm
	token, err := jwt.Parse(tokenInput.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Check whether token is valid
	if !token.Valid {
		h.unauthorizedResponse(w, r)
		return
	}

	// Get token claims and parse it
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		errInvalidToken := errors.New("invalid token")

		// Get refresh UUID and convert the interface to string
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			fmt.Println("uuid")
			h.errorResponse(w, r, http.StatusUnprocessableEntity, errInvalidToken.Error())
			return
		}

		// Get user id from claims and convert it to int
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			fmt.Println("user id")
			h.errorResponse(w, r, http.StatusUnprocessableEntity, errInvalidToken.Error())
			return
		}

		// Remove user token from cache
		deleted, err := h.removeUserTokenFromCache(refreshUuid)
		if err != nil || deleted == 0 {
			h.errorResponse(w, r, http.StatusUnprocessableEntity, errInvalidToken.Error())
			return
		}

		// Create a new pair of tokens
		ts, err := h.createToken(int(userId))
		if err != nil {
			h.errorResponse(w, r, http.StatusForbidden, err.Error())
			return
		}

		// Save new pair of tokens in cache
		err = h.saveTokenInformation(int(userId), ts)
		if err != nil {
			h.errorResponse(w, r, http.StatusForbidden, err.Error())
			return
		}

		// Return new tokens
		tokens := jsonResponse{
			"accessToken":  ts.AccessToken,
			"refreshToken": ts.RefreshToken,
		}

		if err := sendJSON(w, tokens, http.StatusOK, nil); err != nil {
			h.internalErrorResponse(w, r, err)
		}
		return
	}

	h.errorResponse(w, r, http.StatusUnauthorized, "token expired")
}
