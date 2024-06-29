package api

import (
	s "auth/internal/structs"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var DevicesMap sync.Map
var UserCodesMap sync.Map

// DeviceAuthorize adds new device
func DeviceAuthorize(c *gin.Context) {
	var response s.DeviceAuthorizationResponse

	deviceUUID := c.Request.URL.Query().Get("client_id")

	if deviceUUID == "" {
		var errorResponse s.ResponseError
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "client_id is invalid"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	DevicesMap.Range(func(key, value interface{}) bool {
		if key == deviceUUID {
			response = value.(s.DeviceAuthorizationResponse)
			return false
		}
		return true
	})

	if response.DeviceCode != "" {
		c.IndentedJSON(http.StatusOK, response)
		return
	}

	codeFirstPartLen := int(float32(deviceVerificationCodeLength)*0.5 + 0.5)

	response.DeviceCode = uuid.New().String()
	response.UserCode = generateRandomCode(codeFirstPartLen) + "-" + generateRandomCode(deviceVerificationCodeLength-codeFirstPartLen)
	response.VerificationURI = deviceVerificationURI
	response.VerificationURIComplete = deviceVerificationURI + "?user_code=" + response.UserCode
	response.ExpiresIn = deviceVerificationCodeExpiration
	response.Interval = deviceVerificationCodeAttempts

	UserCodesMap.Store(response.UserCode, deviceUUID)
	DevicesMap.Store(deviceUUID, response)

	var devicesMapSize int
	DevicesMap.Range(func(key, value interface{}) bool {
		devicesMapSize++
		return true
	})

	if devicesMapSize == 1 {
		go controlExpiration()
	}

	c.IndentedJSON(http.StatusOK, response)
}

// generateRandomCode generates random code
func generateRandomCode(numSymbols int) (code string) {
	for i := 0; i < numSymbols; i++ {
		code += string(randomLetterOrDigit())
	}

	return
}

// randomLetterOrDigit generates random symbols
func randomLetterOrDigit() rune {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rune(deviceVerificationCodeCharSet[r.Intn(len(deviceVerificationCodeCharSet))])
}

// controlExpiration controls the expiration of the user codes
func controlExpiration() {
	var devicesMapSize int = 1

	for devicesMapSize > 0 {
		time.Sleep(1 * time.Second)

		DevicesMap.Range(func(key, value interface{}) bool {
			if response, ok := value.(s.DeviceAuthorizationResponse); ok {
				if response.ExpiresIn == 0 {
					UserCodesMap.Delete(response.UserCode)
					DevicesMap.Delete(key)
				} else {
					response.ExpiresIn--
					DevicesMap.Store(key, response)
				}
			}
			return true
		})

		DevicesMap.Range(func(key, value interface{}) bool {
			devicesMapSize++
			return true
		})
	}
}
