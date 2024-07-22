package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"profile/internal/db"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

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

	var profile s.Profile
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

// UserChangeEmail changes the user email in auth service and if successful, updates the user profile email in the database.
func UserChangeEmail(c *gin.Context) {
	var errorResponse s.ResponseError

	tmpEmail, exists := c.Get("email")
	if !exists || tmpEmail.(string) == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	var changeEmail struct {
		NewEmail string `json:"new_email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&changeEmail); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid request"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	url := server.AuthServerURL + server.ChangePathEmail
	client := &http.Client{}

	payload, err := json.Marshal(changeEmail)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating change email request"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating change email request"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	accessToken, _ := c.Cookie("access_token")
	request.Header.Set("Cookie", "access_token="+accessToken)

	response, err := client.Do(request)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error sending change email request"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error changing email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if err := db.UserEmailChange(tmpEmail.(string), changeEmail.NewEmail); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error updating user profile email in the database"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("User email changed successfully for an ID: ", tmpEmail.(string))
	loger.Debug("New email: ", changeEmail.NewEmail)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Profile email changed successfully", "email": changeEmail.NewEmail})
}
