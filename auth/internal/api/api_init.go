package api

import (
	"os"
	"strconv"

	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var host string
var port string
var deviceVerificationURI string
var deviceVerificationCodeCharSet string
var deviceVerificationCodeLength int
var deviceVerificationCodeExpiration int
var deviceVerificationCodeAttempts int

// ApiInit starts the web service
func ApiInit() {
	readJWTConfig()

	host = cfg.Config["auth.host"]
	if host == "" {
		panic("auth.host is not set")
	}

	port = cfg.Config["auth.port"]
	if port == "" {
		panic("auth.port is not set")
	}

	deviceVerificationURI = cfg.Config["auth.device.verification.url"]
	if deviceVerificationURI == "" {
		panic("auth.device.verification.url is not set")
	}

	deviceVerificationCodeCharSet = cfg.Config["auth.device.verification.code.charset"]
	if deviceVerificationCodeCharSet == "" {
		panic("auth.device.verification.code.charset is not set")
	}

	var err error

	deviceVerificationCodeLength, err = strconv.Atoi(cfg.Config["auth.device.verification.code.length"])
	if err != nil {
		panic("auth.device.verification.code.length is not set")
	}

	deviceVerificationCodeExpiration, err = strconv.Atoi(cfg.Config["auth.device.verification.code.expiration"])
	if err != nil {
		panic("auth.device.verification.code.expiration is not set")
	}

	deviceVerificationCodeAttempts, err = strconv.Atoi(cfg.Config["auth.device.verification.code.attempts"])
	if err != nil {
		panic("auth.device.verification.code.attempts is not set")
	}

	if os.Getenv("DEBUG") != "1" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/api/v1/auth/user/signup", UserSignUp)
	router.GET("/api/v1/auth/user/identify", UserIdentify)
	router.DELETE("/api/v1/auth/user/delete", UserDelete)
	router.POST("/api/v1/auth/user/login", UserLogIn)
	router.POST("/api/v1/auth/user/refresh", UserRefresh)
	router.POST("/api/v1/auth/user/logout", UserLogOut)
	router.POST("/api/v1/auth/user/password", UserPasswordChange)
	router.POST("/api/v1/auth/user/email", UserEmailChange)
	router.POST("/api/v1/auth/device/authorize", DeviceAuthorize)
	router.POST("/api/v1/auth/device/verify", DeviceVerify)
	router.POST("/api/v1/auth/device/token", DeviceTokens)
	router.GET("/api/v1/auth/device/identify", DeviceIdentify)
	router.DELETE("/api/v1/auth/device/delete", DeviceDelete)
	router.POST("/api/v1/auth/device/refresh", DeviceRefresh)

	url := host + ":" + port
	loger.Success("Service successfully started", url)
	router.Run(url)
}
