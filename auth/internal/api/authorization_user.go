package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Log in user
func LogInUser(c *gin.Context) {
	var logIn LoginUserRequest
	var errorResponse ResponseError

	if err := c.BindJSON(&logIn); err != nil {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid json"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	fmt.Println(logIn.Email, logIn.Password, logIn.Remember)

	if logIn.Email == "" || logIn.Password == "" || len(logIn.Password) < 8 {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid email or password"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	var accessToken string
	var refreshToken string
	var err error

	if logIn.Remember {
		accessToken, err = GenerateToken(logIn.Email, 60)
		if err != nil {
			errorResponse.Error = "token_error"
			errorResponse.ErrorDescription = "Failed to generate access token"
			c.IndentedJSON(http.StatusInternalServerError, errorResponse)
			return
		}

		refreshToken, err = GenerateToken(logIn.Email, 24*60*7)
		if err != nil {
			errorResponse.Error = "token_error"
			errorResponse.ErrorDescription = "Failed to generate refresh token"
			c.IndentedJSON(http.StatusInternalServerError, errorResponse)
			return
		}
	} else {
		accessToken, err = GenerateToken(logIn.Email, 15)
		if err != nil {
			errorResponse.Error = "token_error"
			errorResponse.ErrorDescription = "Failed to generate access token"
			c.IndentedJSON(http.StatusInternalServerError, errorResponse)
			return
		}

		refreshToken, err = GenerateToken(logIn.Email, 24*60)
		if err != nil {
			errorResponse.Error = "token_error"
			errorResponse.ErrorDescription = "Failed to generate refresh token"
			c.IndentedJSON(http.StatusInternalServerError, errorResponse)
			return
		}
	}

	// db.AddNewUser(newUser.Email, newUser.Password, accessToken, refreshToken)

	response := AuthUserResponse{
		Message:      "User logged in successfully",
		Email:        logIn.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// fmt.Println(response)

	c.JSON(http.StatusOK, response)
}
