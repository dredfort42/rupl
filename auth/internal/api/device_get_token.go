package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDeviceAccessToken get device access token
func GetDeviceAccessToken(c *gin.Context) {
	var errorResponse ResponseError

	grantType := c.Request.URL.Query().Get("grant_type")
	deviceCode := c.Request.URL.Query().Get("device_code")
	clientID := c.Request.URL.Query().Get("client_id")

	if grantType == "" || deviceCode == "" || clientID == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing required parameter"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if grantType != "urn:ietf:params:oauth:grant-type:device_code" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid grant type"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if Devices == nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid device code"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	} else if _, ok := Devices[clientID]; !ok {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid client id"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	} else if Devices[clientID].DeviceCode != deviceCode {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid device code"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if DeviceAccessTokens == nil ||
		DeviceAccessTokens[clientID].AccessToken == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Device access token not found"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	c.IndentedJSON(http.StatusOK, DeviceAccessTokens[clientID])

	delete(Devices, clientID)
	delete(DeviceAccessTokens, clientID)
}
