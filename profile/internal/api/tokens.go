package api

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dredfort42/tools/logprinter"
)

// ValidateAccessToken validates access token
func ValidateAccessToken(accesstToken string) bool {
	authServerURL := fmt.Sprintf("http://%s:%s/api/v1/auth/verify", config["auth.host"], config["auth.port"])
	client := &http.Client{}

	request, err := http.NewRequest("GET", authServerURL, nil)
	if err != nil {
		logprinter.PrintError("Error creating request", err)
		return false
	}

	request.Header.Set("Cookie", fmt.Sprintf("access_token=%s", accesstToken))

	response, err := client.Do(request)
	if err != nil {
		logprinter.PrintError("Error sending request", err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return true
	} else {
		return false
	}
}

// ParseToken parses token
func ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		logprinter.PrintError("Error parsing token", err)
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	return userID, nil
}
