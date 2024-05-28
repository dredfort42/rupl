package api

import (
	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ApiInit starts the web service
func ApiInit() {

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/api/v1/auth/register", RegisterUser)
	router.POST("/api/v1/auth/login", LogInUser)
	// router.POST("/api/v1/auth/device_authorization", DeviceAuthorization)
	// router.DELETE("/api/v1/auth/device_authorization", DeviceDeauthorization)
	// router.POST("/api/v1/auth/device_identify", DeviceIdentification)
	// router.POST("/api/v1/auth/device_token", GetDeviceAccessToken)
	// router.GET("/api/v1/auth/device_verify", VerifyDevice)
	router.GET("/api/v1/auth/verify", VerifyUser)
	// router.GET("/api/v1/auth/verify_email", VerifyEmail)
	router.GET("/api/v1/auth/refresh", RefreshUserTokens)
	router.GET("/api/v1/auth/logout", LogOutUser)

	url := cfg.Config["auth.host"] + ":" + cfg.Config["auth.port"]

	loger.Success("Service successfully started", url)
	router.Run(url)
}
