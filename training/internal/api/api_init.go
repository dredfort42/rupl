package api

import (
	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ResponseError is a struct for JSON error
type ResponseError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// ApiInit starts the web service
func ApiInit() {

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/api/v1/training/session", CreateSession)
	// router.GET("/api/v1/profile", GetProfile)
	// router.GET("/api/v1/profile/devices", GetDevices)
	// router.POST("/api/v1/profile/devices", CreateDevice)
	// router.PUT("/api/v1/profile/devices", UpdateDevice)
	// router.DELETE("/api/v1/profile/devices", DeleteDevice)

	router.GET("/api/v1/training/task", GetTask)
	router.POST("/api/v1/training/task", DeclineTask)

	url := cfg.Config["training.host"] + ":" + cfg.Config["training.port"]

	loger.Success("Service successfully started", url)
	router.Run(url)
}
