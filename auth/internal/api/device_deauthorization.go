package api

import (
	"net/http"

	"auth/internal/db"

	"github.com/gin-gonic/gin"
)

func DeviceDeauthorization(c *gin.Context) {
	var clientID string
	var tClientID string
	var accessToken string
	var errorResponse ResponseError
	var err error

	clientID = c.Request.URL.Query().Get("client_id")
	accessToken = c.Request.URL.Query().Get("access_token")

	if clientID == "" || accessToken == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing required parameter"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if tClientID, err = ParseToken(accessToken); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to parse access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if clientID != tClientID {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Invalid access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if !db.CheckDeviceAccessToken(clientID, accessToken) {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Invalid device access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if !db.DeleteDeviceAccessToken(db.GetEmailByAccessToken(accessToken)) {
		errorResponse.Error = "internal_error"
		errorResponse.ErrorDescription = "Failed to deauthorize device"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device successfully deauthorized"})
}
