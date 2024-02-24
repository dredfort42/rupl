package main

import (
	"os"

	"auth/internal/api"
	"auth/internal/db"
	"github.com/dredfort42/tools/configreader"
)

func main() {
	config, err := configreader.GetConfig()

	if err != nil {
		os.Exit(1)
	}

	db.Start(config)
	api.Start(config)
}
