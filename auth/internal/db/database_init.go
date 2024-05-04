package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
)

// Database is the database struct
type Database struct {
	database   *sql.DB
	tableUsers string
	err        error
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

var DEBUG bool = false

// DatabaseInit initializes the database
func DatabaseInit() {
	db.tableUsers = cfg.Config["db.table.users"]

	connectToDatabase()
	checkTables()
}
