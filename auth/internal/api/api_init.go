package api

import (
	"os"

	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ApiInit starts the web service
func ApiInit() {
	if os.Getenv("DEBUG") != "1" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/api/v1/auth/user/signup", UserSignUp)
	router.POST("/api/v1/auth/user/login", UserLogIn)
	router.POST("/api/v1/auth/user/logout", LogOutUser) // remove tockens separately and all together
	router.DELETE("/api/v1/auth/user/delete", UserDelete)
	router.GET("/api/v1/auth/user/refresh", RefreshUserTokens)
	router.GET("/api/v1/auth/user/verify", VerifyUser)
	// router.POST("/api/v1/auth/device_authorization", DeviceAuthorization)
	// router.DELETE("/api/v1/auth/device_authorization", DeviceDeauthorization)
	// router.POST("/api/v1/auth/device_identify", DeviceIdentification)
	// router.POST("/api/v1/auth/device_token", GetDeviceAccessToken)
	// router.GET("/api/v1/auth/device_verify", VerifyDevice)
	// router.GET("/api/v1/auth/verify_email", VerifyEmail)

	url := cfg.Config["auth.host"] + ":" + cfg.Config["auth.port"]

	loger.Success("Service successfully started", url)
	router.Run(url)
}
