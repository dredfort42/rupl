package api

import (
	"errors"
	"net/http"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
)

type ClientType uint8

const (
	_ ClientType = iota
	USER
	DEVICE
)

// AuthMiddleware checks for a token in the request header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var clientType ClientType

		userToken, _ := c.Cookie("access_token")
		deviceToken := c.GetHeader("Authorization")
		clientID := c.Request.URL.Query().Get("client_id")

		if userToken != "" && deviceToken == "" {
			clientType = USER
		} else if userToken == "" && deviceToken != "" {
			clientType = DEVICE
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var email string

		switch clientType {
		case USER:
			email, err := UserIdentify(userToken)
			if email == "" || err != nil {
				c.Abort()
				return
			}
		case DEVICE:
			if clientID == "" {
				c.Abort()
				return
			}

			email, err := DeviceIdentify(clientID, deviceToken)
			if email == "" || err != nil {
				c.Abort()
				return
			}
		}

		c.Set("email", email)
		c.Next()
	}
}

// UserIdentify verifies the user based on the access token provided in the request.
func UserIdentify(accessToken string) (string, error) {
	url := server.AuthServerURL + server.IdentifyPathUser
	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		loger.Error("Error creating request", err)
		return "", err
	}

	request.Header.Set("Cookie", "access_token="+accessToken)

	response, err := client.Do(request)
	if err != nil {
		loger.Error("Error sending request", err)
		return "", err
	}
	defer response.Body.Close()

	var email string
	if response.StatusCode == http.StatusOK {
		email = response.Header.Get("email")
		if email == "" {
			return "", errors.New("email not found")
		}
	}

	return email, nil
}

// DeviceIdentify verifies the device based on the client ID and access token provided in the request.
func DeviceIdentify(clientID string, accessToken string) (string, error) {
	url := server.AuthServerURL + server.IdentifyPathDevice + "?client_id=" + clientID
	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		loger.Error("Error creating request", err)
		return "", err
	}

	request.Header.Set("Authorization", accessToken)

	response, err := client.Do(request)
	if err != nil {
		loger.Error("Error sending request", err)
		return "", err
	}
	defer response.Body.Close()

	var email string
	if response.StatusCode == http.StatusOK {
		email = response.Header.Get("email")
		if email == "" {
			return "", errors.New("email not found")
		}
	}

	return email, nil
}
