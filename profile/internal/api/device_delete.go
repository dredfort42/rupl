package api

import (
	"net/http"
	"profile/internal/db"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

// DeviceDelete deletes a device based on the access token provided in the request.
func DeviceDelete(c *gin.Context) {
	var device s.Device
	var errorResponse s.ResponseError

	tmpEmail, exists := c.Get("email")
	if !exists || tmpEmail.(string) == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	err := c.ShouldBindJSON(&device)
	if err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid request"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if device.DeviceUUID == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing device UUID"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	exists = db.DeviceExistsCheck(tmpEmail.(string), device.DeviceUUID)
	if !exists {
		errorResponse.Error = "not_found"
		errorResponse.ErrorDescription = "Device not found"
		c.IndentedJSON(http.StatusNotFound, errorResponse)
		return
	}

	if err = db.DeviceDelete(tmpEmail.(string), device.DeviceUUID); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error deleting user device"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("Device deleted successfully for an ID: ", device.DeviceUUID)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}
