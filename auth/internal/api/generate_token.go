package api

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken generates token
func GenerateToken(userID string, minitesToExpire time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)
	tokenClaims["user_id"] = userID
	tokenClaims["exp"] = time.Now().Add(time.Minute * minitesToExpire).Unix()
	response, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return response, nil
}

// // GenerateTokens generates refresh token
// func GenerateRefreshToken(userID string, hoursToExpire time.Duration) (string, error) {
// 	refreshToken := jwt.New(jwt.SigningMethodHS256)
// 	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
// 	refreshClaims["user_id"] = userID
// 	refreshClaims["exp"] = time.Now().Add(time.Hour * hoursToExpire).Unix()
// 	refresh, err := refreshToken.SignedString([]byte("secret"))
// 	if err != nil {
// 		return "", err
// 	}

// 	return refresh, nil
// }
