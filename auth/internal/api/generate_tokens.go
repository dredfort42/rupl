package api

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateTokens generates access token
func GenerateAccessToken(userID string) (string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["user_id"] = userID
	accessClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	access, err := accessToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return access, nil
}

// GenerateTokens generates refresh token
func GenerateRefreshToken(userID string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["user_id"] = userID
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	refresh, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return refresh, nil
}
