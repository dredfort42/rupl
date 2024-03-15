package api

import (
	"net/http"

	// "profile/internal/db"

	"github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

// Get user profile
func GetProfile(c *gin.Context) {
	var accessToken string
	var errorResponse ResponseError
	var err error

	if accessToken, err = c.Cookie("access_token"); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if !ValidateAccessToken(accessToken) {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "User not authenticated"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	email, err := ParseToken(accessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to parse access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if DEBUG {
		logprinter.PrintInfo("Request a user profile for an ID: ", email)
	}

	// access_token, ok := c.")

	// if !ok {
	// 	errorResponse.Error = "invalid_request"
	// 	errorResponse.ErrorDescription = "User not authenticated"
	// 	c.IndentedJSON(http.StatusUnauthorized, errorResponse)
	// 	return
	// }

	// userProfile, err := db.GetProfile(userID.(int))
	// if err != nil {
	// 	errorResponse.Error = "server_error"
	// 	errorResponse.ErrorDescription = "Error getting user profile"
	// 	c.IndentedJSON(http.StatusInternalServerError, errorResponse)
	// 	return
	// }

	// profile, err := db.GetProfile(email)
	// if err != nil {
	// 	c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "profile_error", "error_description": "Failed to get profile from the database"})
	// 	return
	// }

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Profile retrieved successfully", "profile": gin.H{"email": email}})
}