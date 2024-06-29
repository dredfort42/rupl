package api

import (
	s "auth/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserIdentify identifies user
func UserIdentify(c *gin.Context) {
	var accessToken string
	var errorResponse s.ResponseError
	var err error

	accessToken, err = c.Cookie("access_token")
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	_, err = verifyToken(accessToken, s.AccessToken)
	if err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "failed to verify access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user successfully identified"})
}
