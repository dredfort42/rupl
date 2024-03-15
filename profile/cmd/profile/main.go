package main

import (
	"os"

	"profile/internal/api"
	"profile/internal/db"

	"github.com/dredfort42/tools/configreader"
	"github.com/dredfort42/tools/logprinter"
)

func main() {
	if debug := os.Getenv("DEBUG"); debug != "" {
		logprinter.PrintWarning("Debugging is enabled.", "")
	}

	config, err := configreader.GetConfig()
	if err != nil {
		os.Exit(1)
	}

	db.Start(config)
	api.Start(config)
}
