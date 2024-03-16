package api

import (
	"net/http"

	"profile/internal/db"

	"github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

// GetProfile retrieves the user profile based on the access token provided in the request.
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

	userProfile, err := db.GetProfile(email)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error getting user profile"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	logprinter.PrintInfo("User profile retrieved successfully for an ID: ", email)
	logprinter.PrintInfo("User profile: ", userProfile.FirstName+" "+userProfile.LastName)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Profile retrieved successfully", "profile": gin.H{"email": email}})
}

// CreateProfile creates a new user profile based on the access token provided in the request.
func CreateProfile(c *gin.Context) {
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
		logprinter.PrintInfo("Request to create a user profile for an ID: ", email)
	}

	var profile Profile
	if err := c.BindJSON(&profile); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid request"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	var profileDB db.Profile
	profileDB.Email = email
	profileDB.FirstName = profile.FirstName
	profileDB.LastName = profile.LastName
	profileDB.DateOfBirth = profile.DateOfBirth
	profileDB.Gender = profile.Gender

	if err := db.CreateProfile(profileDB); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating user profile"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	logprinter.PrintInfo("User profile created successfully for an ID: ", email)
	logprinter.PrintInfo("User profile: ", profile.FirstName+" "+profile.LastName)

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Profile created successfully", "profile": gin.H{"email": email}})
}
