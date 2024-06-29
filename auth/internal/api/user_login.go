package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserLogIn logs in a user
func UserLogIn(c *gin.Context) {
	var logIn s.LoginUserRequest
	var accessToken string
	var refreshToken string
	var errorResponse s.ResponseError
	var err error

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "authorization header is invalid"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "authorization header is invalid"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	payload, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "authorization header is invalid"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "authorization header is invalid"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	logIn.Email, logIn.Password = pair[0], pair[1]

	if !db.IsUserPasswordCorrect(logIn.Email, logIn.Password) {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "invalid email or password"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	remember := c.Query("remember")
	if remember == "true" {
		logIn.Remember = true
	}

	if logIn.Remember {
		accessToken, refreshToken, err = getTokens(logIn.Email, s.JWTConfig.BrowserAccessTokenExpiration, s.JWTConfig.BrowserRefreshTokenExpiration)
		c.SetCookie("access_token", accessToken, s.JWTConfig.BrowserAccessTokenExpiration, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, s.JWTConfig.BrowserRefreshTokenExpiration, "/", "", false, true)
	} else {
		accessToken, refreshToken, err = getTokens(logIn.Email, s.JWTConfig.OneTimeAccessTokenExpiration, s.JWTConfig.OneTimeRefreshTokenExpiration)
		c.SetCookie("access_token", accessToken, s.JWTConfig.OneTimeAccessTokenExpiration, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, s.JWTConfig.OneTimeRefreshTokenExpiration, "/", "", false, true)
	}

	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to generate tokens | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if logIn.Remember {
		err = db.SessionCreate(logIn.Email, accessToken, refreshToken, false)
	} else {
		err = db.SessionUpdateOneTime(logIn.Email, accessToken, refreshToken)
	}

	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to update user tokens in the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user successfully logged in"})
}
