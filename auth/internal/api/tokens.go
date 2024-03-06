package api

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken generates token
func GenerateToken(userID string, minitesToExpire int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)
	tokenClaims["user_id"] = userID
	tokenClaims["exp"] = time.Now().Add(time.Minute * time.Duration(minitesToExpire)).Unix()
	response, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return response, nil
}

// GetAccessAndRefreshTokens gets access and refresh tokens
func GetAccessAndRefreshTokens(email string, accessTokenMinitesToExpire int, refreshTokenMinitesToExpire int) (string, string, error) {
	accessToken, err := GenerateToken(email, accessTokenMinitesToExpire)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateToken(email, refreshTokenMinitesToExpire)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseToken parses token
func ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	return userID, nil
}

// TokenHasExpired checks if token has expired
func TokenHasExpired(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return false, err
	}

	claims := token.Claims.(jwt.MapClaims)
	exp := claims["exp"].(float64)
	if time.Now().Unix() > int64(exp) {
		return true, nil
	}

	return false, nil
}

// RefreshTokens with new access and refresh tokens
func RefreshTokens(refreshToken string, accessTokenMinitesToExpire int, refreshTokenMinitesToExpire int) (string, string, error) {
	userID, err := ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	tokenHasExpired, err := TokenHasExpired(refreshToken)
	if err != nil {
		return "", "", err
	} else if tokenHasExpired {
		return "", "", errors.New("Refresh token has expired")
	}

	return GetAccessAndRefreshTokens(userID, accessTokenMinitesToExpire, refreshTokenMinitesToExpire)
}
