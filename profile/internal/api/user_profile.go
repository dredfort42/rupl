package api

import (
	"net/http"

	// "auth/internal/db"

	"github.com/gin-gonic/gin"
)

// Get user profile
func GetUserProfile(c *gin.Context) {
	// var accessToken string
	// var errorResponse ResponseError
	// var err error

	// if accessToken, err = c.Cookie("access_token"); err != nil {
	// 	errorResponse.Error = "token_error"
	// 	errorResponse.ErrorDescription = "Missing access token"
	// 	c.IndentedJSON(http.StatusUnauthorized, errorResponse)
	// 	return
	// }

	// access_token, ok := c.")

	// if !ok {
	// 	errorResponse.Error = "invalid_request"
	// 	errorResponse.ErrorDescription = "User not authenticated"
	// 	c.IndentedJSON(http.StatusUnauthorized, errorResponse)
	// 	return
	// }

	// userProfile, err := db.GetUserProfile(userID.(int))
	// if err != nil {
	// 	errorResponse.Error = "server_error"
	// 	errorResponse.ErrorDescription = "Error getting user profile"
	// 	c.IndentedJSON(http.StatusInternalServerError, errorResponse)
	// 	return
	// }

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User profile"})
}
