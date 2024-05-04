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
	tableProfiles string
	tableDevices  string
	err           error
}

var db Database

// connectToDatabase connects to the database and returns a pointer to it
func connectToDatabase() {
	url := "host=" + cfg.Config["db.host"] +
		" port=" + cfg.Config["db.port"] +
		" user=" + cfg.Config["db.user"] +
		" password=" + cfg.Config["db.password"] +
		" dbname=" + cfg.Config["db.database.name"] +
		" sslmode=" + cfg.Config["db.security.ssl"]

	db.database, db.err = sql.Open("postgres", url)

	if db.err != nil {
		db.database.Close()
		panic(db.err)
	} else if db.err = db.database.Ping(); db.err != nil {
		db.database.Close()
		panic(db.err)
	}

	loger.Success("Successfully connected to database")
}

// DatabaseInit initializes the database
func DatabaseInit() {
	db.tableProfiles = cfg.Config["db.table.profiles"]
	db.tableDevices = cfg.Config["db.table.devices"]

	connectToDatabase()
	checkTables()
}
