package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserRefresh refreshes user tokens
func UserRefresh(c *gin.Context) {
	var email string
	var refreshToken string
	var newAccessToken string
	var newRefreshToken string
	var errorResponse s.ResponseError
	var err error

	refreshToken, err = c.Cookie("refresh_token")
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "missing refresh token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	email, err = verifyToken(refreshToken, s.RefreshToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to verify refresh token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if db.IsOneTimeRefreshToken(refreshToken) {
		newAccessToken, newRefreshToken, err = getTokens(email, 15*60, 24*60*60)
		c.SetCookie("access_token", newAccessToken, 15*60, "/", "", false, true)
		c.SetCookie("refresh_token", newRefreshToken, 24*60*60, "/", "", false, true)
	} else {
		newAccessToken, newRefreshToken, err = getTokens(email, 60*60, 24*60*7*60)
		c.SetCookie("access_token", newAccessToken, 60*60, "/", "", false, true)
		c.SetCookie("refresh_token", newRefreshToken, 24*60*7*60, "/", "", false, true)
	}

	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to generate tokens | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	err = db.SessionUpdate(email, refreshToken, newAccessToken, newRefreshToken)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to update user tokens in the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user tokens successfully refreshed"})
}
