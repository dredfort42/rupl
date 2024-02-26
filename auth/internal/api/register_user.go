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

	if err := c.BindJSON(&newUser); err != nil {
		var errorResponse ResponseError

		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid JSON"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if newUser.Email == "" || newUser.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter email or password"})
		return
	}

	db.AddNewUser(newUser.Email, newUser.Password)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "User successfully registered"})
}
