package api

import (
	"net/http"
	"profile/internal/db"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

// UserGet retrieves the user profile based on the access token provided in the request.
func UserGet(c *gin.Context) {
	var errorResponse s.ResponseError

	tmpEmail, exists := c.Get("email")
	if !exists || tmpEmail.(string) == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	userProfile, err := db.UserGet(tmpEmail.(string))
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error getting user profile"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if userProfile.Email == "" {
		errorResponse.Error = "not_found"
		errorResponse.ErrorDescription = "Profile not found"
		c.IndentedJSON(http.StatusNotFound, errorResponse)
		return
	}

	loger.Debug("User profile retrieved successfully for an ID: ", userProfile.Email)
	loger.Debug("User profile: ", userProfile.FirstName+" "+userProfile.LastName+" "+userProfile.DateOfBirth+" "+userProfile.Gender)

	c.IndentedJSON(http.StatusOK, userProfile)
}
