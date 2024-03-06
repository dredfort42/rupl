package api

import (
	"net/http"

	"auth/internal/db"

	"github.com/gin-gonic/gin"
)

// Verify user
func VerifyUser(c *gin.Context) {
	var email string
	var accessToken string
	var errorResponse ResponseError
	var err error

	// Get token from cookie or header
	if accessToken, err = c.Cookie("access_token"); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if email, err = ParseToken(accessToken); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to parse access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if !db.CheckUserAccessToken(email, accessToken) {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Invalid access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully verified"})
}
