package api

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/dredfort42/tools/configreader"
	"github.com/dredfort42/tools/logprinter"
)

var config configreader.ConfigMap

// Start starts the web service
func Start(configMap configreader.ConfigMap) {
	config = configMap

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/api/v1/auth/register", RegisterUser)
	router.POST("/api/v1/auth/login", LogInUser)
	router.POST("/api/v1/auth/device_authorization", DeviceAuthorization)
	router.GET("/api/v1/auth/verify", VerifyUser)
	router.GET("/api/v1/auth/refresh", RefreshUserTokens)
	router.GET("/api/v1/auth/logout", LogOutUser)

	url := fmt.Sprintf("%s:%s", config["auth.host"], config["auth.port"])

	logprinter.PrintSuccess("Entry point", url)
	router.Run(url)
}
