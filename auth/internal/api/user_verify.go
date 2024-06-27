package api

import (
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Verify user
func UserVerify(c *gin.Context) {
	var accessToken string
	var errorResponse s.ResponseError
	var err error

	accessToken, err = c.Cookie("access_token")
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	_, err = verifyToken(accessToken, s.AccessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to verify access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user successfully verified"})
}

// // Verify email
// func VerifyEmail(c *gin.Context) {
// 	var email string
// 	var accessToken string
// 	var errorResponse ResponseError
// 	var err error

// 	if accessToken, err = c.Cookie("access_token"); err != nil {
// 		errorResponse.Error = "token_error"
// 		errorResponse.ErrorDescription = "Missing access token"
// 		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
// 		return
// 	}

// 	if email, err = ParseToken(accessToken); err != nil {
// 		errorResponse.Error = "token_error"
// 		errorResponse.ErrorDescription = "Failed to parse access token"
// 		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
// 		return
// 	}

// 	if !db.CheckUserEmailVerified(email) {
// 		errorResponse.Error = "email_error"
// 		errorResponse.ErrorDescription = "Not verified email"
// 		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Email verified"})
// }
