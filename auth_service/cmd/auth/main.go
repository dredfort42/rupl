package main

import (
	"os"

	"auth/internal/api"
	"github.com/dredfort42/tools/configreader"
	// "github.com/dredfort42/tools/logprinter"
)

func main() {
	config, err := configreader.GetConfig()

	if err != nil {
		os.Exit(1)
	}

	api.StartService(config)
}
