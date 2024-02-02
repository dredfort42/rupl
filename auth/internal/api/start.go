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
	router.POST("/api/v1/auth/device_authorization", DeviceAuthorization)
	// router.GET("/cameras", GetAllEntries)
	// router.GET("/cameras/:uuid", GetEntryByUUID)
	// router.POST("/cameras/update", UpdateEntry)
	// router.DELETE("/cameras/:uuid", DeleteEntryByUUID)

	url := fmt.Sprintf("%s:%s", config["auth.host"], config["auth.port"])

	logprinter.PrintSuccess("Entry point", url)
	router.Run(url)

}
