package api

import (
	"fmt"
	"os"

	"github.com/dredfort42/tools/configreader"
	"github.com/dredfort42/tools/logprinter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var config configreader.ConfigMap

var DEBUG bool = false

// Start starts the web service
func Start(configMap configreader.ConfigMap) {
	if debug := os.Getenv("DEBUG"); debug != "" {
		DEBUG = true
	}

	config = configMap

	router := gin.Default()
	router.Use(cors.Default())
	// router.POST("/api/v1/auth/register", RegisterUser)
	// router.POST("/api/v1/auth/login", LogInUser)
	// router.POST("/api/v1/auth/device_authorization", DeviceAuthorization)
	router.GET("/api/v1/profile", GetProfile)
	// router.GET("/api/v1/auth/refresh", RefreshUserTokens)
	// router.GET("/api/v1/auth/logout", LogOutUser)

	url := fmt.Sprintf("%s:%s", config["profile.host"], config["profile.port"])

	logprinter.PrintSuccess("Entry point", url)
	router.Run(url)
}
