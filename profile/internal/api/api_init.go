package api

import (
	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ApiInit starts the profile service
func ApiInit() {

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/api/v1/profile", GetProfile)
	router.POST("/api/v1/profile", CreateProfile)
	router.GET("/api/v1/profile/devices", GetDevices)
	router.POST("/api/v1/profile/devices", CreateDevice)
	router.PUT("/api/v1/profile/devices", UpdateDevice)
	router.DELETE("/api/v1/profile/devices", DeleteDevice)

	url := cfg.Config["profile.host"] + ":" + cfg.Config["profile.port"]

	loger.Success("Service successfully started", url)
	router.Run(url)
}
