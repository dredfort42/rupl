package db

import (
	"database/sql"

	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
	_ "github.com/lib/pq"
)

// Database is the database struct
type Database struct {
	database      *sql.DB
	tableUsers    string
	tableSessions string
	tableDevices  string
}

var db Database

// DatabaseInit initializes the database
func DatabaseInit() {
	db.tableUsers = cfg.Config["db.table.users"]
	if db.tableUsers == "" {
		panic("Table users is not set")
	}

	db.tableSessions = cfg.Config["db.table.sessions"]
	if db.tableSessions == "" {
		panic("Table sessions is not set")
	}

	db.tableDevices = cfg.Config["db.table.devices"]
	if db.tableDevices == "" {
		panic("Table devices is not set")
	}

	databaseConnect()
	tablesCheck()
	databaseCleanerStart()

	loger.Success("Database successfully initialized")
}
