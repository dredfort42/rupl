package api

import (
	db "auth/internal/db"
	s "auth/internal/structs"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// getToken generates token
// id: id
// expired: token expiration time in seconds
func getToken(id string, expiration int) (response string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)
	tokenClaims["user_id"] = id
	tokenClaims["exp"] = time.Now().Add(time.Second * time.Duration(expiration)).Unix()

	return token.SignedString([]byte("secret"))
}

// getTokens gets access and refresh tokens
// id: id
// accessTokenExpiration: access token expiration time in seconds
// refreshTokenExpiration: refresh token expiration time in seconds
func getTokens(id string, accessTokenExpiration int, refreshTokenExpiration int) (accessToken string, refreshToken string, err error) {
	accessToken, err = getToken(id, accessTokenExpiration)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = getToken(id, refreshTokenExpiration)
	if err != nil {
		return "", "", err
	}

	return
}

// parseToken verifies token
func parseToken(token string) (id string, err error) {
	var jwtToken *jwt.Token

	jwtToken, err = jwt.Parse(
		token,
		func(jwtToken *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
	if err != nil {
		err = errors.New("failed to parse token")
		return
	}

	if time.Now().Unix() > int64(jwtToken.Claims.(jwt.MapClaims)["exp"].(float64)) {
		err = errors.New("token has expired")
		return
	}

	id = jwtToken.Claims.(jwt.MapClaims)["user_id"].(string)
	return
}

// verifyToken verifies token
func verifyToken(token string, tokenType s.TokenType) (id string, err error) {
	id, err = parseToken(token)
	if err != nil {
		return
	}

	db.IsTokenExists(id, token, tokenType)

	return
}

// // isTokenExpired checks if token has expired
// func isTokenExpired(tokenString string) (result bool, err error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return []byte("secret"), nil
// 	})
// 	if err != nil {
// 		return
// 	}

// 	result = time.Now().Unix() > int64(token.Claims.(jwt.MapClaims)["exp"].(float64))

// 	return
// }

// // refreshTokens with new access and refresh tokens
// func refreshTokens(refreshToken string, accessTokenMinitesToExpire int, refreshTokenMinitesToExpire int) (newAccessToken string, newRefreshToken string, err error) {
// 	var userID string

// 	userID, err = parseToken(refreshToken)
// 	if err != nil {
// 		return
// 	}

// 	tokenHasExpired, err := isTokenExpired(refreshToken)
// 	if err != nil {
// 		return
// 	} else if tokenHasExpired {
// 		err = errors.New("refresh token has expired")

// 		return
// 	}

// 	return GetAccessAndRefreshTokens(userID, accessTokenMinitesToExpire, refreshTokenMinitesToExpire)
// }
