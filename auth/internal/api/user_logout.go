package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Log out user
func LogOutUser(c *gin.Context) {
	var email string
	var accessToken string
	var errorResponse s.ResponseError
	var err error

	accessToken, err = c.Cookie("access_token")
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	email, err = verifyToken(accessToken, s.AccessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = err.Error()
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	err = db.DeleteTokens(email)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to delete user's tokens from the database"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User successfully logged out"})
}
