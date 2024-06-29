package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var DeviceTokensMap sync.Map

// DeviceVerify make the device identify and add it to the users devices
func DeviceVerify(c *gin.Context) {
	var email string
	var accessToken string
	var errorResponse s.ResponseError
	var err error

	accessToken, err = c.Cookie("access_token")
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	userCode := c.Request.URL.Query().Get("user_code")
	if userCode == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "user_code is required"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	email, err = verifyToken(accessToken, s.AccessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to verify access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	deviceUUID, ok := UserCodesMap.Load(userCode)
	if !ok {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "user_code is invalid"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	UserCodesMap.Delete(userCode)

	deviceAccessToken, deviceRefreshToken, err := getTokens(deviceUUID.(string), s.JWTConfig.DeviceAccessTokenExpiration, s.JWTConfig.DeviceRefreshTokenExpiration)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to generate device tokens | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	err = db.DeviceCreate(email, deviceUUID.(string), deviceAccessToken, deviceRefreshToken)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to create a new device in the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	DeviceTokensMap.Store(deviceUUID, s.DeviceAccessTokenResponse{
		TokenType:    "Bearer",
		AccessToken:  deviceAccessToken,
		ExpiresIn:    s.JWTConfig.DeviceAccessTokenExpiration,
		RefreshToken: deviceRefreshToken,
	})

	c.IndentedJSON(http.StatusOK, gin.H{"message": "device successfully verified"})
}
