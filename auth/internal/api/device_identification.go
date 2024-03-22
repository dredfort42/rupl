package api

import (
	"net/http"

	"auth/internal/db"

	"github.com/gin-gonic/gin"
)

var DeviceAccessTokens = make(map[string]DeviceAccessTokenResponse)

// DeviceIdentification make the device identify and add it to the users devices
func DeviceIdentification(c *gin.Context) {
	var email string
	var accessToken string
	var errorResponse ResponseError
	var err error

	if accessToken, err = c.Cookie("access_token"); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if email, err = ParseToken(accessToken); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to parse access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	tokenHasExpired, err := TokenHasExpired(accessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to check access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	} else if tokenHasExpired {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Access token has expired"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if !db.CheckUserAccessToken(email, accessToken) {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Invalid access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	userCode := c.Request.URL.Query().Get("user_code")

	if userCode == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing required parameter: user_code"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	clientID, ok := UserCodes[userCode]
	if !ok {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid user code"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	deviceAccessToken, err := GenerateToken(clientID, 60*24*365)
	if err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Failed to generate device access token"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	db.UpdateUserDeviceUUID(email, clientID)
	db.UpdateUserDeviceAccessToken(email, deviceAccessToken)

	DeviceAccessTokens[clientID] = DeviceAccessTokenResponse{
		AccessToken: deviceAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   60 * 24 * 365,
	}

	delete(UserCodes, userCode)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Verification completed successfully"})
}
