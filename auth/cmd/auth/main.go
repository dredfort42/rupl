package main

import (
	"auth/internal/api"
	"auth/internal/db"

	cfg "github.com/dredfort42/tools/configreader"
)

func main() {
	err := cfg.GetConfig()
	if err != nil {
		panic(err)
	}

	db.DatabaseInit()
	api.ApiInit()
}
