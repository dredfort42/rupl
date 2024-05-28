package api

import (
	"auth/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Log out user
func LogOutUser(c *gin.Context) {
	var email string
	var accessToken string
	var errorResponse ResponseError
	var result bool
	var err error

	accessToken, err = c.Cookie("access_token")
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)

		return
	}

	email, err = ParseToken(accessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to parse access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)

		return
	}

	result, err = db.IsUserAccessTokenExists(email, accessToken)
	if !result {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Invalid access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)

		return
	}

	err = db.DeleteUserTokens(email)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "Failed to delete user's tokens from the database"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		return
	}

	// Delete tokens from cookie
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User successfully logged out"})
}
