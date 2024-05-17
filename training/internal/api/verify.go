package api

import (
	"net/http"

	s "training/internal/structs"

	"github.com/gin-gonic/gin"
)

// // VerifyUser verifies the user based on the access token provided in the request.
// func VerifyUser(c *gin.Context) string {
// 	var accessToken string
// 	var errorResponse s.ResponseError
// 	var err error

// 	if accessToken, err = c.Cookie("access_token"); err != nil {
// 		errorResponse.Error = "token_error"
// 		errorResponse.ErrorDescription = "Missing access token"
// 		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
// 		return ""
// 	}

// 	email := ValidateAccessToken(accessToken)
// 	if email == "" {
// 		errorResponse.Error = "invalid_request"
// 		errorResponse.ErrorDescription = "User not authenticated"
// 		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
// 		return ""
// 	}

// 	return email
// }

// VerifyDevice verifies the device based on the client ID and access token provided in the request.
func VerifyDevice(c *gin.Context) string {
	var accessToken string
	var errorResponse s.ResponseError

	if accessToken = c.Request.URL.Query().Get("access_token"); accessToken == "" {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return ""
	}

	email := ValidateAccessToken(accessToken)
	if email == "" {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "User not authenticated"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return ""
	}

	return email
}
