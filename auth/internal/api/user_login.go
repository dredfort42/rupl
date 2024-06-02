package api

import (
	"auth/internal/db"
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserLogIn logs in a user
func UserLogIn(c *gin.Context) {
	var logIn s.LoginUserRequest
	var accessToken string
	var refreshToken string
	var errorResponse s.ResponseError
	var err error

	err = c.BindJSON(&logIn)
	if err != nil {
		errorResponse.Error = "invalid_json"
		errorResponse.ErrorDescription = "invalid JSON in the request body"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if logIn.Email == "" || logIn.Password == "" || len(logIn.Password) < 8 ||
		!db.IsUserExists(logIn.Email) || !db.IsPasswordCorrect(logIn.Email, logIn.Password) {
		errorResponse.Error = "invalid_parameter"
		errorResponse.ErrorDescription = "invalid email or password"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if logIn.Remember {
		accessToken, refreshToken, err = getTokens(logIn.Email, 60*60, 24*60*7*60)
		c.SetCookie("access_token", accessToken, 60*60, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, 24*60*7*60, "/", "", false, true)
	} else {
		accessToken, refreshToken, err = getTokens(logIn.Email, 15*60, 24*60*60)
		c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "", false, true)
	}

	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to generate tokens"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if logIn.Remember {
		err = db.RememberUserTokens(logIn.Email, accessToken, refreshToken)
	} else {
		err = db.UpdateUserTokens(logIn.Email, accessToken, refreshToken)
	}

	if err != nil {
		errorResponse.Error = "database_error"
		errorResponse.ErrorDescription = "failed to update user tokens in the database"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user successfully logged in"})
}
