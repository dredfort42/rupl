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

	var profile Profile
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

// UserUpdate updates the user profile based on the access token provided in the request.
func UserUpdate(c *gin.Context) {
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

	var profile Profile
	if err := c.BindJSON(&profile); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid request"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if profile.FirstName == "" {
		profile.FirstName = userProfile.FirstName
	}

	if profile.LastName == "" {
		profile.LastName = userProfile.LastName
	}

	if profile.DateOfBirth == "" {
		profile.DateOfBirth = userProfile.DateOfBirth
	}

	if profile.Gender == "" {
		profile.Gender = userProfile.Gender
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

	if err := db.UserUpdate(profileDB); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error updating user profile"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("User profile updated successfully for an ID: ", profileDB.Email)
	loger.Debug("User profile: ", profile.FirstName+" "+profile.LastName+" "+profile.DateOfBirth+" "+profile.Gender)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "profile": profileDB})
}
