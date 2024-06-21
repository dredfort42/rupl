package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RefreshTokens refreshes user tokens
func RefreshUserTokens(c *gin.Context) {
	var email string
	var accessToken string
	var refreshToken string
	var errorResponse s.ResponseError
	var result bool
	var err error

	refreshToken, err = c.Cookie("refresh_token")
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "missing refresh token | " + err.Error()
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	email, err = verifyToken(accessToken, s.RefreshToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to verify refresh token " + err.Error()
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	var isRefreshTokenRemembered bool

	isRefreshTokenRemember, err = db.

		// result, err = db.IsUserRememberMeSet(email)
		// if result {
		// 	accessToken, refreshToken, err = GetAccessAndRefreshTokens(email, 60, 24*60*7)
		// 	c.SetCookie("access_token", accessToken, 60*60, "/", "", false, true)
		// 	c.SetCookie("refresh_token", refreshToken, 24*60*60*7, "/", "", false, true)
		// } else {
		// 	accessToken, refreshToken, err = GetAccessAndRefreshTokens(email, 15, 24*60)
		// 	c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)
		// 	c.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "", false, true)
		// }

		// if err != nil {
		// 	errorResponse.Error = "token_error"
		// 	errorResponse.ErrorDescription = "Failed to generate tokens"
		// 	c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		// 	return
		// }

		// err = db.UpdateUserTokens(email, accessToken, refreshToken)
		// if err != nil {
		// 	errorResponse.Error = "database_error"
		// 	errorResponse.ErrorDescription = "Failed to update user tokens"
		// 	c.IndentedJSON(http.StatusInternalServerError, errorResponse)

		// 	return
		// }

		c.JSON(http.StatusOK, gin.H{"message": "User tokens successfully refreshed"})
}
