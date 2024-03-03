package api

import (
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

// // RefreshToken refreshes token
// func RefreshToken(tokenString string) (string, error) {
// 	userID, err := ParseToken(tokenString)
// 	if err != nil {
// 		return "", err
// 	}

// 	accessToken, err := GenerateToken(userID, 60)
// 	if err != nil {
// 		return "", err
// 	}

// 	return accessToken, nil
// }
