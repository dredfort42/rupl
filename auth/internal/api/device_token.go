package api

import (
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeviceTokens get device access token
func DeviceTokens(c *gin.Context) {
	var errorResponse s.ResponseError

	grantType := c.Request.URL.Query().Get("grant_type")
	deviceCode := c.Request.URL.Query().Get("device_code")
	deviceUUID := c.Request.URL.Query().Get("client_id")

	if grantType == "" || deviceCode == "" || deviceUUID == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "missing required parameter"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if grantType != "urn:ietf:params:oauth:grant-type:device_code" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "grant_type is invalid"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	DevicesMap.Range(func(key, value interface{}) bool {
		if key == deviceUUID {
			if value.(s.DeviceAuthorizationResponse).DeviceCode != deviceCode {
				errorResponse.Error = "invalid_request"
				errorResponse.ErrorDescription = "device_code is invalid"
				c.IndentedJSON(http.StatusBadRequest, errorResponse)
				return false
			}
			return false
		}
		return true
	})

	DeviceTokensMap.Range(func(key, value interface{}) bool {
		if key == deviceUUID {
			c.IndentedJSON(http.StatusOK, value)
			return false
		}
		return true
	})

	DevicesMap.Delete(deviceUUID)
	DeviceTokensMap.Delete(deviceUUID)
}
