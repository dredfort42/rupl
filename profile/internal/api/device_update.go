package api

import (
	"net/http"
	"profile/internal/db"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

// DeviceUpdate updates a device based on the access token provided in the request.
func DeviceUpdate(c *gin.Context) {
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

	err = device.DeviceStructCheck()
	if err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Device struct check failed | " + err.Error()
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

	err = db.DeviceUpdate(tmpEmail.(string), device)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error updating user device"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("Device updated successfully for an ID: ", device.DeviceUUID)
	loger.Debug("Device name: ", device.DeviceName)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Device updated successfully", "device": device})
}
