package api

import (
	"net/http"
	"profile/internal/db"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

const (
	GENDER_MAN   = "man"
	GENDER_WOMAN = "woman"
	GENDER_OTHER = "other"
)

// UserCreate creates a new user profile based on the access token provided in the request.
func UserCreate(c *gin.Context) {
	var errorResponse s.ResponseError

	tmpEmail, exists := c.Get("email")
	if !exists || tmpEmail.(string) == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	var profile s.Profile
	if err := c.BindJSON(&profile); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid request"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if profile.FirstName == "" || profile.LastName == "" || profile.DateOfBirth == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing required fields"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if profile.Gender != GENDER_MAN && profile.Gender != GENDER_WOMAN {
		profile.Gender = GENDER_OTHER
	}

	var profileDB s.User
	profileDB.Email = tmpEmail.(string)
	profileDB.FirstName = profile.FirstName
	profileDB.LastName = profile.LastName
	profileDB.DateOfBirth = profile.DateOfBirth
	profileDB.Gender = profile.Gender

	if err := db.UserExistsCheck(profileDB.Email); err {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Profile already exists"
		c.IndentedJSON(http.StatusConflict, errorResponse)
		return
	}

	if err := db.UserCreate(profileDB); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating user profile"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("User profile created successfully for an ID: ", profileDB.Email)
	loger.Debug("User profile: ", profile.FirstName+" "+profile.LastName+" "+profile.DateOfBirth+" "+profile.Gender)

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Profile created successfully", "profile": profileDB})
}
