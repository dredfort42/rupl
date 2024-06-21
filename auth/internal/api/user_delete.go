package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserDelete deletes a user
func UserDelete(c *gin.Context) {
	var email string
	var accessToken string
	var errorResponse s.ResponseError
	var err error

	accessToken, err = c.Cookie("access_token")
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "missing access token | " + err.Error()
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	email, err = verifyToken(accessToken, s.AccessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to verify access token | " + err.Error()
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	err = db.DeleteUser(email)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to delete user from the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.Header("email", email)
	c.JSON(http.StatusOK, gin.H{"message": "user successfully deleted"})
}
