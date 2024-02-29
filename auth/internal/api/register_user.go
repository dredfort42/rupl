package api

import (
	"fmt"

	"net/http"

	"auth/internal/db"
	"github.com/gin-gonic/gin"
)

// RegisterUser adds new user
func RegisterUser(c *gin.Context) {
	var newUser RegisterUserRequest
	var errorResponse ResponseError

	if err := c.BindJSON(&newUser); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid json"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if newUser.Email == "" {
		errorResponse.Error = "missing_required_parameter"
		errorResponse.ErrorDescription = "email"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if newUser.Password == "" {
		errorResponse.Error = "missing_required_parameter"
		errorResponse.ErrorDescription = "password"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if len(newUser.Password) < 8 {
		errorResponse.Error = "password_error"
		errorResponse.ErrorDescription = "Password must be at least 8 characters long"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	db.AddNewUser(newUser.Email, newUser.Password)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "User successfully registered"})
}
