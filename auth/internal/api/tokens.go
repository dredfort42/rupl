package api

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken generates token
func GenerateToken(userID string, minitesToExpire int) (response string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)
	tokenClaims["user_id"] = userID
	tokenClaims["exp"] = time.Now().Add(time.Minute * time.Duration(minitesToExpire)).Unix()

	return token.SignedString([]byte("secret"))
}

// GetAccessAndRefreshTokens gets access and refresh tokens
func GetAccessAndRefreshTokens(email string, accessTokenMinitesToExpire int, refreshTokenMinitesToExpire int) (accessToken string, refreshToken string, err error) {
	accessToken, err = GenerateToken(email, accessTokenMinitesToExpire)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = GenerateToken(email, refreshTokenMinitesToExpire)
	if err != nil {
		return "", "", err
	}

	return
}

// ParseToken parses token
func ParseToken(tokenString string) (userID string, err error) {
	var token *jwt.Token

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}

	userID = token.Claims.(jwt.MapClaims)["user_id"].(string)

	return
}

// TokenHasExpired checks if token has expired
func TokenHasExpired(tokenString string) (result bool, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}

	result = time.Now().Unix() > int64(token.Claims.(jwt.MapClaims)["exp"].(float64))

	return
}

// RefreshTokens with new access and refresh tokens
func RefreshTokens(refreshToken string, accessTokenMinitesToExpire int, refreshTokenMinitesToExpire int) (newAccessToken string, newRefreshToken string, err error) {
	var userID string

	userID, err = ParseToken(refreshToken)
	if err != nil {
		return
	}

	tokenHasExpired, err := TokenHasExpired(refreshToken)
	if err != nil {
		return
	} else if tokenHasExpired {
		err = errors.New("refresh token has expired")

		return
	}

	return GetAccessAndRefreshTokens(userID, accessTokenMinitesToExpire, refreshTokenMinitesToExpire)
}
