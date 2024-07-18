package db

import (
	cfg "github.com/dredfort42/tools/configreader"
	_ "github.com/lib/pq"
)

// DatabaseInit initializes the database
func DatabaseInit() {
	db.tableUsers = cfg.Config["db.table.profile.users"]
	if db.tableUsers == "" {
		panic("db.table.profile.users is empty")
	}

	db.tableDevices = cfg.Config["db.table.profile.devices"]
	if db.tableDevices == "" {
		panic("db.table.profile.devices is empty")
	}

	databaseConnect()
	tablesCheck()
}
