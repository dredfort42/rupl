package api

import (
	"net/http"
	"profile/internal/db"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

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

	devices, err := db.DevicesGet(tmpEmail.(string))
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error getting user devices"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("User devices retrieved successfully for an ID: ", devices.Email)

	c.IndentedJSON(http.StatusOK, devices)
}
