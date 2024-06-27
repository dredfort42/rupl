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
	readJWTConfig()

	host := cfg.Config["auth.host"]
	if host == "" {
		panic("auth.host is not set")
	}

	port := cfg.Config["auth.port"]
	if port == "" {
		panic("auth.port is not set")
	}

	if os.Getenv("DEBUG") != "1" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/api/v1/auth/user/signup", UserSignUp)
	router.GET("/api/v1/auth/user/verify", UserVerify)
	router.DELETE("/api/v1/auth/user/delete", UserDelete)
	router.POST("/api/v1/auth/user/login", UserLogIn)
	router.GET("/api/v1/auth/user/refresh", UserRefresh)
	router.POST("/api/v1/auth/user/logout", UserLogOut)
	// router.POST("/api/v1/auth/device_authorization", DeviceAuthorization)
	// router.DELETE("/api/v1/auth/device_authorization", DeviceDeauthorization)
	// router.POST("/api/v1/auth/device_identify", DeviceIdentification)
	// router.POST("/api/v1/auth/device_token", GetDeviceAccessToken)
	// router.GET("/api/v1/auth/device_verify", VerifyDevice)
	// router.GET("/api/v1/auth/verify_email", VerifyEmail)

	url := host + ":" + port
	loger.Success("Service successfully started", url)
	router.Run(url)
}
