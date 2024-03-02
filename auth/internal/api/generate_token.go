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
