package main

import (
	"os"
	"time"

	"github.com/dredfort42/tools/configreader"
	"github.com/dredfort42/tools/logprinter"
)

func main() {
	config, err := configreader.GetConfig()

	if err != nil {
		os.Exit(1)
	}

	_ = config

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		logprinter.PrintInfo("Profile Manager is started", "")
		<-ticker.C
	}
}
