package api

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"manager/internal/configuration"
)

var config configuration.ConfigMap

// Start starts the web service
func Start(configMap configuration.ConfigMap) {

	config = configMap

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/cameras", GetAllEntries)
	router.GET("/cameras/:uuid", GetEntryByUUID)
	router.POST("/cameras/add", AddNewEntry)
	router.POST("/cameras/update", UpdateEntry)
	router.DELETE("/cameras/:uuid", DeleteEntryByUUID)

	url := fmt.Sprintf("%s:%s", config["manager.host"], config["manager.port"])

	router.Run(url)
}
