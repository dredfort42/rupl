package api

import (
	"net/http"
	"profile/internal/db"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

// ProfileGet retrieves the user profile based on the access token provided in the request.
func ProfileGet(c *gin.Context) {
	var errorResponse s.ResponseError

	email := c.Request.URL.Query().Get("email")
	if email == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	userProfile, err := db.ProfileGet(email)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error getting user profile"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("User profile retrieved successfully for an ID: ", userProfile.Email)
	loger.Debug("User profile: ", userProfile.FirstName+" "+userProfile.LastName)

	c.IndentedJSON(http.StatusOK, userProfile)
}

// CreateProfile creates a new user profile based on the access token provided in the request.
func CreateProfile(c *gin.Context) {
	var errorResponse s.ResponseError

	email := c.Request.URL.Query().Get("email")
	if email == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("Request to create a user profile for an ID: ", email)

	var profile Profile
	if err := c.BindJSON(&profile); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid request"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	var profileDB s.Profile
	profileDB.Email = email
	profileDB.FirstName = profile.FirstName
	profileDB.LastName = profile.LastName
	profileDB.DateOfBirth = profile.DateOfBirth
	profileDB.Gender = profile.Gender

	if err := db.ProfileCreate(profileDB); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating user profile"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("User profile created successfully for an ID: ", email)
	loger.Debug("User profile: ", profile.FirstName+" "+profile.LastName)

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Profile created successfully", "profile": profileDB})
}
