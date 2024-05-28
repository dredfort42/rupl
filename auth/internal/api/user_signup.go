package api

import (
	"auth/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUser adds new user
func RegisterUser(c *gin.Context) {
	var newUser UserCredentials
	var errorResponse ResponseError
	var accessToken string
	var refreshToken string
	var result bool
	var err error

	err = c.BindJSON(&newUser)
	if err != nil {
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

	result, err = db.IsUserExists(newUser.Email)
	if result {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "User with this email already exists"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		return
	}

	accessToken, refreshToken, err = GetAccessAndRefreshTokens(newUser.Email, 15, 24*60)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to generate tokens"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		return
	}

	err = db.AddNewUser(newUser.Email, newUser.Password, accessToken, refreshToken)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "Failed to add new user to the database"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		return
	}

	c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User successfully registered"})
}
