package model

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
)

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type TokenMetadata struct {
	AccessUuid string
	UserId     int
}

func (a *Auth) CreateToken(userId int) (*TokenDetails, error) {
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
	tokenExpTimeMinutes := time.Duration(viper.GetInt("server.auth.tokenExpTime"))
	refreshExpTimeDays := time.Duration(viper.GetInt("server.auth.refreshExpTime"))

	// Token details structure
	td := TokenDetails{}

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
	accessClaims["refresh_uuid"] = td.AccessUuid
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
