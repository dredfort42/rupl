package api

import (
	"auth/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Verify user
func VerifyUser(c *gin.Context) {
	var userID string
	var accessToken string
	var errorResponse ResponseError
	var err error

	if accessToken, err = c.Cookie("access_token"); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if userID, err = ParseToken(accessToken); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to parse access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	tokenHasExpired, err := TokenHasExpired(accessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to check access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	} else if tokenHasExpired {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Access token has expired"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if !db.CheckUserAccessToken(userID, accessToken) &&
		!db.CheckDeviceAccessToken(userID, accessToken) {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Invalid access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	c.Header("email", db.GetEmailByAccessToken(accessToken))

	c.JSON(http.StatusOK, gin.H{"message": "User successfully verified"})
}

// Verify email
func VerifyEmail(c *gin.Context) {
	var email string
	var accessToken string
	var errorResponse ResponseError
	var err error

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

	if !db.CheckUserEmailVerified(email) {
		errorResponse.Error = "email_error"
		errorResponse.ErrorDescription = "Not verified email"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified"})
}
