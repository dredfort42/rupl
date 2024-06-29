package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserDelete deletes a user
func UserDelete(c *gin.Context) {
	var deleteRequest s.UserCredentials
	var email string
	var errorResponse s.ResponseError

	err := c.BindJSON(&deleteRequest)
	if err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "invalid JSON in the request body | " + err.Error()
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

	email, err = verifyToken(accessToken, s.AccessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to verify access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	if deleteRequest.Email == "" || deleteRequest.Password == "" || deleteRequest.Email != email {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "email or password is invalid"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if !db.IsUserPasswordCorrect(deleteRequest.Email, deleteRequest.Password) {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "password is invalid"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	err = db.SessionDeleteAll(email)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to delete user's sessions from the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	err = db.DeviceDeleteAll(email)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to delete user's devices from the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	err = db.UserDelete(email)
	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to delete user from the database | " + err.Error()
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user successfully deleted"})
}
