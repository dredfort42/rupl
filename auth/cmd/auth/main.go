package main

import (
	"os"

	"auth/internal/api"
	"github.com/dredfort42/tools/configreader"
)

func main() {
	config, err := configreader.GetConfig()

	if err != nil {
		os.Exit(1)
	}

	api.Start(config)
}
