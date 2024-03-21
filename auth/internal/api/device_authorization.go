package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var Devices = make(map[string]DeviceAuthorizationResponse)
var UserCodes = make(map[string]string)

// DeviceAuthorization adds new device
func DeviceAuthorization(c *gin.Context) {
	var response DeviceAuthorizationResponse

	clientID := c.Request.URL.Query().Get("client_id")

	if clientID == "" || len(clientID) < 32 || len(clientID) > 36 {
		var errorResponse ResponseError

		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing required parameter: client_id"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if Devices != nil {
		if response, ok := Devices[clientID]; ok {
			c.IndentedJSON(http.StatusOK, response)
			return
		}
		fmt.Println("Devices check not ok")
	}

	url := fmt.Sprintf("%s://%s", config["entrypoint.protocol.ssl"], config["entrypoint.address"])

	response.DeviceCode = uuid.New().String()
	response.UserCode = fmt.Sprintf("%04d-%04d", generateRandomDigits(4), generateRandomDigits(4))
	response.VerificationURI = url + "/device"
	response.VerificationURIComplete = url + "/device?user_code=" + response.UserCode
	response.ExpiresIn = 600
	response.Interval = 5

	UserCodes[response.UserCode] = clientID
	Devices[clientID] = response

	if len(Devices) == 1 {
		go controlExpiration()
	}

	c.IndentedJSON(http.StatusOK, response)
}

func generateRandomDigits(numDigits int) int {
	maxValue := int64(pow(10, numDigits) - 1)
	randomNumber := rand.Int63n(maxValue + 1)

	return int(randomNumber)
}

func pow(base, exponent int) int {
	result := 1
	for i := 0; i < exponent; i++ {
		result *= base
	}
	return result
}

func controlExpiration() {
	for len(Devices) > 0 {
		time.Sleep(1 * time.Second)
		for key, value := range Devices {
			if value.ExpiresIn == 0 {
				delete(UserCodes, value.UserCode)
				delete(Devices, key)
			} else {
				value.ExpiresIn--
				Devices[key] = value
			}
		}
	}
}
