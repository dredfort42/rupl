package api

import (
	"fmt"
	"os"

	"github.com/dredfort42/tools/configreader"
	"github.com/dredfort42/tools/logprinter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var DEBUG bool = false
var config configreader.ConfigMap

// Start starts the web service
func Start(configMap configreader.ConfigMap) {
	if debug := os.Getenv("DEBUG"); debug != "" {
		DEBUG = true
	}

	config = configMap

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/api/v1/training/task", GetTask)
	router.POST("/api/v1/training/task", DeclineTask)
	// router.GET("/api/v1/profile", GetProfile)
	// router.POST("/api/v1/profile", CreateProfile)
	// router.GET("/api/v1/profile/devices", GetDevices)
	// router.POST("/api/v1/profile/devices", CreateDevice)
	// router.PUT("/api/v1/profile/devices", UpdateDevice)
	// router.DELETE("/api/v1/profile/devices", DeleteDevice)

	url := fmt.Sprintf("%s:%s", config["training.host"], config["training.port"])

	logprinter.PrintSuccess("Entry point", url)
	router.Run(url)
}
