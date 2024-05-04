package api

import (
	"auth/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RefreshTokens refreshes user tokens
func RefreshUserTokens(c *gin.Context) {
	var email string
	var accessToken string
	var refreshToken string
	var errorResponse ResponseError
	var err error

	if refreshToken, err = c.Cookie("refresh_token"); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Missing refresh token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if email, err = ParseToken(refreshToken); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to parse refresh token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	tokenHasExpired, err := TokenHasExpired(refreshToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to check refresh token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	} else if tokenHasExpired {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Refresh token has expired"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if !db.CheckUserRefreshToken(email, refreshToken) {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Invalid refresh token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if db.CheckUserRememberMe(email) {
		accessToken, refreshToken, err = GetAccessAndRefreshTokens(email, 60, 24*60*7)
		c.SetCookie("access_token", accessToken, 60*60, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, 24*60*60*7, "/", "", false, true)
	} else {
		accessToken, refreshToken, err = GetAccessAndRefreshTokens(email, 15, 24*60)
		c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "", false, true)
	}

	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to generate tokens"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	db.UpdateUserTokens(email, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{"message": "Tokens successfully refreshed"})
}
