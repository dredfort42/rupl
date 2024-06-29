package api

import (
	db "auth/internal/db"
	s "auth/internal/structs"
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	cfg "github.com/dredfort42/tools/configreader"
)

// readJWTConfig reads JWT configuration
func readJWTConfig() {
	s.JWTConfig = s.JWTParamethers{
		TokenSecret: cfg.Config["jwt.secret"],
	}
	if s.JWTConfig.TokenSecret == "" {
		panic("JWT secret is not set")
	}

	var expiration int
	var err error

	expiration, err = strconv.Atoi(cfg.Config["jwt.onetime.access.token.expiration"])
	if err != nil {
		panic("JWT onetime access token expiration is not set")
	}
	s.JWTConfig.OneTimeAccessTokenExpiration = expiration

	expiration, err = strconv.Atoi(cfg.Config["jwt.onetime.refresh.token.expiration"])
	if err != nil {
		panic("JWT onetime refresh token expiration is not set")
	}
	s.JWTConfig.OneTimeRefreshTokenExpiration = expiration

	expiration, err = strconv.Atoi(cfg.Config["jwt.browser.access.token.expiration"])
	if err != nil {
		panic("JWT browser access token expiration is not set")
	}
	s.JWTConfig.BrowserAccessTokenExpiration = expiration

	expiration, err = strconv.Atoi(cfg.Config["jwt.browser.refresh.token.expiration"])
	if err != nil {
		panic("JWT browser refresh token expiration is not set")
	}
	s.JWTConfig.BrowserRefreshTokenExpiration = expiration

	expiration, err = strconv.Atoi(cfg.Config["jwt.device.access.token.expiration"])
	if err != nil {
		panic("JWT device access token expiration is not set")
	}
	s.JWTConfig.DeviceAccessTokenExpiration = expiration

	expiration, err = strconv.Atoi(cfg.Config["jwt.device.refresh.token.expiration"])
	if err != nil {
		panic("JWT device refresh token expiration is not set")
	}
	s.JWTConfig.DeviceRefreshTokenExpiration = expiration
}

// getToken generates token
func getToken(id string, expiration int) (response string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)
	tokenClaims["id"] = id
	tokenClaims["exp"] = time.Now().Add(time.Second * time.Duration(expiration)).Unix()

	return token.SignedString([]byte(s.JWTConfig.TokenSecret))
}

// getTokens gets access and refresh tokens
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
			return []byte(s.JWTConfig.TokenSecret), nil
		})
	if err != nil {
		err = errors.New("failed to parse token")
		return
	}

	if time.Now().Unix() > int64(jwtToken.Claims.(jwt.MapClaims)["exp"].(float64)) {
		err = errors.New("token has expired")
		return
	}

	id = jwtToken.Claims.(jwt.MapClaims)["id"].(string)
	if id == "" {
		err = errors.New("failed to get user ID from token")
	}

	return
}

// verifyToken verifies token
func verifyToken(token string, tokenType s.TokenType) (id string, err error) {
	id, err = parseToken(token)
	if err != nil {
		return
	}

	if !db.IsTokenExist(id, token, tokenType) {
		err = errors.New("token does not exist")
	}

	return
}
