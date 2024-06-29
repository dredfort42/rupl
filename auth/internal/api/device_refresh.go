package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeviceRefresh refreshes the device's access and refresh tokens
func DeviceRefresh(c *gin.Context) {
	var deviceUUID string
	var errorResponse s.ResponseError
	var err error

	clientID := c.Request.URL.Query().Get("client_id")
	grantType := c.Request.URL.Query().Get("grant_type")
	refreshToken := c.Request.URL.Query().Get("refresh_token")

	if clientID == "" || grantType == "" || refreshToken == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "missing required parameter"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if grantType != "refresh_token" {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "grant_type is invalid"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	deviceUUID, err = verifyToken(refreshToken, s.DeviceRefreshToken)
	if err != nil || deviceUUID != clientID {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to verify device refresh token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	deviceAccessToken, deviceRefreshToken, err := getTokens(deviceUUID, s.JWTConfig.DeviceAccessTokenExpiration, s.JWTConfig.DeviceRefreshTokenExpiration)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to generate device tokens | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	err = db.DeviceRefresh(deviceUUID, deviceAccessToken, deviceRefreshToken)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to update device tokens in the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	response := s.DeviceAccessTokenResponse{
		TokenType:    "Bearer",
		AccessToken:  deviceAccessToken,
		ExpiresIn:    s.JWTConfig.DeviceAccessTokenExpiration,
		RefreshToken: deviceRefreshToken,
	}

	c.IndentedJSON(http.StatusOK, response)
}
