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

	if !deviceStructCheck(device, c) {
		return
	}

	if err = db.DeviceCreate(email, device); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating device"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("Device created successfully for an ID: ", device.DeviceID)
	loger.Debug("Device name: ", device.DeviceName)

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Device created successfully", "device": device})
}

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

	if !deviceStructCheck(device, c) {
		return
	}

	err = db.DeviceUpdate(device)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error updating device"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("Device updated successfully for an ID: ", device.DeviceID)
	loger.Debug("Device name: ", device.DeviceName)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Device updated successfully", "device": device})
}

func deviceStructCheck(device s.Device, c *gin.Context) (success bool) {
	var errorResponse s.ResponseError

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
	return true
}

// DevicesGet returns all devices associated with the user.
func DevicesGet(c *gin.Context) {
	var errorResponse s.ResponseError

	tmpEmail, exists := c.Get("email")
	if !exists || tmpEmail.(string) == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	email := tmpEmail.(string)

	devices, err := db.DevicesGet(email)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error getting user devices"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("User devices retrieved successfully for an ID: ", devices.Email)

	c.IndentedJSON(http.StatusOK, devices)
}

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

	if device.DeviceID == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing device ID"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if err = db.DeviceDelete(device.DeviceID); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error deleting device"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("Device deleted successfully for an ID: ", device.DeviceID)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Device deleted successfully", "device": device})
}
