package api

import (
	// "fmt"
	"net/http"

	"auth/internal/db"

	"github.com/gin-gonic/gin"
)

// Log in user
func LogInUser(c *gin.Context) {
	var logIn LoginUserRequest
	var errorResponse ResponseError
	var accessToken string
	var refreshToken string
	var err error

	if err := c.BindJSON(&logIn); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid json"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	// fmt.Println(logIn.Email, logIn.Password, logIn.Remember)

	if logIn.Email == "" || logIn.Password == "" || len(logIn.Password) < 8 {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid email or password"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	// Check if user exists
	if !db.CheckUserExists(logIn.Email) {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "User does not exist"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	// Check if password is correct
	if !db.CheckUserPassword(logIn.Email, logIn.Password) {
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

	db.UpdateUserTokens(logIn.Email, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}
