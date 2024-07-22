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

// UserDelete deletes a user profile based on the access token provided in the request.
func UserDelete(c *gin.Context) {
	var errorResponse s.ResponseError

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid request"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	url := server.AuthServerURL + server.DeletePathUser
	client := &http.Client{}

	payload, err := json.Marshal(body)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating delete user request"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	request, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payload))
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating delete user request"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	accessToken, _ := c.Cookie("access_token")
	request.Header.Set("Cookie", "access_token="+accessToken)

	response, err := client.Do(request)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error sending delete user request"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error deleting user"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	err = db.UserDelete(body.Email)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error deleting user"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("User profile deleted successfully for an ID: ", body.Email)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Profile deleted"})
}
