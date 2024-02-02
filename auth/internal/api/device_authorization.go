package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeviceAuthorization adds new device
func DeviceAuthorization(c *gin.Context) {
	var request DeviceAuthorizationRequest
	var response DeviceAuthorizationResponse

	c.BindJSON(&request)

	if request.ClientID == "" || len(request.ClientID) < 32 || len(request.ClientID) > 36 {
		var errorResponse DeviceAuthorizationError

		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Missing required parameter: client_id"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	url := fmt.Sprintf("%s://%s", config["entrypoint.protocol.ssl"], config["entrypoint.address"])

	response.DeviceCode = uuid.New().String()
	response.UserCode = strconv.Itoa(generateRandomDigits(4)) + "-" + strconv.Itoa(generateRandomDigits(4))
	response.VerificationURI = url + "/device"
	response.VerificationURIComplete = url + "/device?user_code=" + response.UserCode
	response.ExpiresIn = 1800
	response.Interval = 5

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
