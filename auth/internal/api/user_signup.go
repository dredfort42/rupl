package api

import (
	"net/http"

	"auth/internal/db"

	"github.com/gin-gonic/gin"
)

// RegisterUser adds new user
func RegisterUser(c *gin.Context) {
	var newUser UserCredentials
	var errorResponse ResponseError
	var accessToken string
	var refreshToken string
	var err error

	if err = c.BindJSON(&newUser); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid json"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if newUser.Email == "" {
		errorResponse.Error = "missing_required_parameter"
		errorResponse.ErrorDescription = "email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if newUser.Password == "" {
		errorResponse.Error = "missing_required_parameter"
		errorResponse.ErrorDescription = "password"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if len(newUser.Password) < 8 {
		errorResponse.Error = "password_error"
		errorResponse.ErrorDescription = "Password must be at least 8 characters long"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if db.CheckUserExists(newUser.Email) {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "User with this email already exists"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if accessToken, refreshToken, err = GetAccessAndRefreshTokens(newUser.Email, 15, 24*60); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to generate tokens"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "", false, true)

	db.AddNewUser(newUser.Email, newUser.Password, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{"message": "User successfully registered"})
}
