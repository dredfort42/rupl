package api

import (
	"net/http"

	"auth/internal/db"

	"github.com/gin-gonic/gin"
)

// RegisterUser adds new user
func RegisterUser(c *gin.Context) {
	var newUser UserCredentials
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

	accessToken, err := GenerateToken(newUser.Email, 15)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to generate access token"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	refreshToken, err := GenerateToken(newUser.Email, 24*60)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Failed to generate refresh token"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	db.AddNewUser(newUser.Email, newUser.Password, accessToken, refreshToken)

	response := AuthUserResponse{
		Message:      "User registered successfully",
		Email:        newUser.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, response)
}
