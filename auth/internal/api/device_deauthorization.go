package api

import (
	"net/http"

	"auth/internal/db"

	"github.com/gin-gonic/gin"
)

func DeviceDeauthorization(c *gin.Context) {
	var errorResponse ResponseError

	clientID := c.Request.URL.Query().Get("client_id")
	accessToken := c.Request.URL.Query().Get("access_token")

	if !db.DeleteDeviceAccessToken(clientID, accessToken) {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Invalid device access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device successfully deauthorized"})
}
