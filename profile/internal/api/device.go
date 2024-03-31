package api

import (
	"net/http"

	"profile/internal/db"

	"github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

// CreateDevice creates a new device based on the access token provided in the request.
func CreateDevice(c *gin.Context) {
	var email string
	var device db.Device
	var errorResponse ResponseError
	var err error

	if email = VerifyDevice(c); email == "" {
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

	if err = db.CreateDevice(email, device); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating device"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if DEBUG {
		logprinter.PrintInfo("Device created successfully for an ID: ", device.DeviceID)
		logprinter.PrintInfo("Device name: ", device.DeviceName)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Device created successfully", "device": device})
}

// GetDevices returns all devices associated with the user.
func GetDevices(c *gin.Context) {
	var devices db.UserDevices
	var errorResponse ResponseError
	var err error

	if clientID := c.Request.URL.Query().Get("client_id"); clientID != "" {
		devices, err = db.GetDevices(VerifyDevice(c))
	} else {
		devices, err = db.GetDevices(VerifyUser(c))
	}

	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error getting user devices"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if DEBUG {
		logprinter.PrintInfo("User devices retrieved successfully for an ID: ", devices.Email)
	}

	c.IndentedJSON(http.StatusOK, devices)
}

// DeleteDevice deletes a device based on the access token provided in the request.
func DeleteDevice(c *gin.Context) {
	var email string
	var device db.Device
	var errorResponse ResponseError
	var err error

	if email = VerifyDevice(c); email == "" {
		return
	}

	if err = c.ShouldBindJSON(&device); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid request"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if device.DeviceID == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing device ID"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if err = db.DeleteDevice(email, device.DeviceID); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error deleting device"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if DEBUG {
		logprinter.PrintInfo("Device deleted successfully for an ID: ", device.DeviceID)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Device deleted successfully", "device": device})
}
