package api

import (
	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ApiInit starts the training service
func ApiInit() {

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/api/v1/session", CreateSession)

	url := cfg.Config["sessions.receiver.host"] + ":" + cfg.Config["sessions.receiver.port"]

	loger.Success("Service successfully started", url)
	router.Run(url)
}
