package api

import (
	"fmt"
	"net/http"

	// "auth/internal/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

// Get user profile
func GetUserProfile(c *gin.Context) {

	var accessToken string
	var errorResponse ResponseError
	var err error

	if accessToken, err = c.Cookie("access_token"); err != nil {
		errorResponse.Error = "token_error"
		errorResponse.ErrorDescription = "Missing access token"
		c.IndentedJSON(http.StatusUnauthorized, errorResponse)
		return
	}

	authServerURL := fmt.Sprintf("http://%s:%s/api/v1/auth/verify", config["auth.host"], config["auth.port"])
	client := &http.Client{}

	req, err := http.NewRequest("GET", authServerURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Cookie", fmt.Sprintf("access_token=%s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		logprinter.PrintSuccess("User authenticated: ", accessToken)
	} else {
		logprinter.PrintError("User not authenticated", nil)
		return
	}

	userID, err := ParseToken(accessToken)

	logprinter.PrintSuccess("User ID: ", userID)

	// access_token, ok := c.")

	// if !ok {
	// 	errorResponse.Error = "invalid_request"
	// 	errorResponse.ErrorDescription = "User not authenticated"
	// 	c.IndentedJSON(http.StatusUnauthorized, errorResponse)
	// 	return
	// }

	// userProfile, err := db.GetUserProfile(userID.(int))
	// if err != nil {
	// 	errorResponse.Error = "server_error"
	// 	errorResponse.ErrorDescription = "Error getting user profile"
	// 	c.IndentedJSON(http.StatusInternalServerError, errorResponse)
	// 	return
	// }

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User profile"})
}

// ParseToken parses token
func ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	return userID, nil
}
