package api

import (
	"auth/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Log in user
func LogInUser(c *gin.Context) {
	var logIn LoginUserRequest
	var errorResponse ResponseError
	var accessToken string
	var refreshToken string
	var result bool
	var err error

	err = c.BindJSON(&logIn)
	if err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid json"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)

		return
	}

	if logIn.Email == "" || logIn.Password == "" || len(logIn.Password) < 8 {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid email or password"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		return
	}

	result, err = db.IsUserExists(logIn.Email)
	if !result {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "User does not exist"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		return
	}

	result, err = db.IsUserPasswordCorrect(logIn.Email, logIn.Password)
	if !result {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid email or password"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		return
	}

	if logIn.Remember {
		accessToken, refreshToken, err = GetAccessAndRefreshTokens(logIn.Email, 60, 24*60*7)
		c.SetCookie("access_token", accessToken, 60*60, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, 24*60*60*7, "/", "", false, true)
	} else {
		accessToken, refreshToken, err = GetAccessAndRefreshTokens(logIn.Email, 15, 24*60)
		c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "", false, true)
	}

	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to generate tokens"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		return
	}

	err = db.UpdateUserRememberMe(logIn.Email, logIn.Remember)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "Failed to update remember_me status"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		return
	}

	err = db.UpdateUserTokens(logIn.Email, accessToken, refreshToken)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "Failed to update tokens"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}
