package api

import (
	"net/http"

	"profile/internal/db"

	"github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

// CreateDevice creates a new device based on the access token provided in the request.
func CreateDevice(c *gin.Context) {
	var device db.Device
	var email string
	var errorResponse ResponseError
	var err error

	email = VerifyDevice(c)
	if email == "" {
		return
	}

	if err = c.ShouldBindJSON(&device); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid request"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if device.DeviceModel == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing device model"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if device.DeviceName == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing device name"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if device.SystemName == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing system name"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if device.SystemVersion == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing system version"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if device.DeviceID == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing device ID"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	// if err = db.CreateDevice(device); err != nil {
	// 	errorResponse.Error = "server_error"
	// 	errorResponse.ErrorDescription = "Error creating device"
	// 	c.IndentedJSON(http.StatusInternalServerError, errorResponse)
	// 	return
	// }

	if DEBUG {
		logprinter.PrintInfo("Device created successfully for an ID: ", device.DeviceID)
		logprinter.PrintInfo("Device name: ", device.DeviceName)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Device created successfully", "profile": gin.H{"email": email}})
}
