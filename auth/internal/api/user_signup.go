package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserSignUp adds new user
func UserSignUp(c *gin.Context) {
	var newUser s.UserCredentials
	var errorResponse s.ResponseError
	var accessToken string
	var refreshToken string
	var result bool
	var err error

	err = c.BindJSON(&newUser)
	if err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "invalid JSON in the request body | " + err.Error()
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if newUser.Email == "" {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "email address is required "
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if !IsUserEmailValid(newUser.Email) {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "email address is invalid"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if newUser.Password == "" {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "password is required"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if IsUserPasswordContainsForbiddenCharacters(newUser.Password) {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "password contains forbidden characters"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if !IsUserPasswordStrength(newUser.Password) {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "password isn't strong enough"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	result = db.IsUserExists(newUser.Email)
	if result {
		errorResponse.Error = "user_exists"
		errorResponse.ErrorDescription = "user with this email already exists"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	accessToken, refreshToken, err = getTokens(newUser.Email, s.JWTConfig.OneTimeAccessTokenExpiration, s.JWTConfig.OneTimeRefreshTokenExpiration)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to generate tokens | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	err = db.UserSignUp(newUser.Email, newUser.Password, accessToken, refreshToken)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "error creating user in the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.SetCookie("access_token", accessToken, s.JWTConfig.OneTimeAccessTokenExpiration, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, s.JWTConfig.OneTimeRefreshTokenExpiration, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "user successfully registered"})
}
