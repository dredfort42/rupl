package api

import (
	"net/http"
	"profile/internal/db"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

// DeviceCreate creates a new device based on the access token provided in the request.
func DeviceCreate(c *gin.Context) {
	var device s.Device
	var errorResponse s.ResponseError

	tmpEmail, exists := c.Get("email")
	if !exists || tmpEmail.(string) == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	email := tmpEmail.(string)

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

	if err = db.DeviceCreate(email, device); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating device"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("Device created successfully for an ID: ", device.DeviceUUID)
	loger.Debug("Device name: ", device.DeviceName)

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Device created successfully", "device": device})
}
