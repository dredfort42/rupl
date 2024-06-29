package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// IsUserEmailValid checks the email address
func IsUserEmailValid(email string) (result bool) {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// UserEmailChange changes the user email address
func UserEmailChange(c *gin.Context) {
	var changeEmailRequest s.UserChangeEmailRequest
	var errorResponse s.ResponseError

	err := c.BindJSON(&changeEmailRequest)
	if err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "invalid JSON in the request body | " + err.Error()
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if changeEmailRequest.NewEmail == "" {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "new email address is required"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if !IsUserEmailValid(changeEmailRequest.NewEmail) {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "new email address is invalid"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	accessToken, err := c.Cookie("access_token")
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	email, err := verifyToken(accessToken, s.AccessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to verify access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if changeEmailRequest.NewEmail == email {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "new email address is the same as the old one"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if changeEmailRequest.Password == "" {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "password is required"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if !db.IsUserPasswordCorrect(email, changeEmailRequest.Password) {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "password is invalid"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	err = db.UserEmailChange(email, changeEmailRequest.NewEmail)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to update user email in the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	err = db.SessionDeleteAll(email)

	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to delete user's tokens from the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "user email successfully changed"})
}
