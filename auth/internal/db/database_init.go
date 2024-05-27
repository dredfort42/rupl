package db

import (
	"database/sql"

	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
	_ "github.com/lib/pq"
)

// Database is the database struct
type Database struct {
	database   *sql.DB
	tableUsers string
}

var db Database

// connectToDatabase connects to the database and returns a pointer to it
func connectToDatabase() (err error) {
	url := "host=" + cfg.Config["db.host"] +
		" port=" + cfg.Config["db.port"] +
		" user=" + cfg.Config["db.user"] +
		" password=" + cfg.Config["db.password"] +
		" dbname=" + cfg.Config["db.database.name"] +
		" sslmode=" + cfg.Config["db.security.ssl"]

	db.database, err = sql.Open("postgres", url)
	if err != nil {
		db.database.Close()
		return
	}

	err = db.database.Ping()
	if err != nil {
		db.database.Close()
		return
	}

	loger.Success("Successfully connected to database")

	return
}

// DatabaseInit initializes the database
func DatabaseInit() (err error) {
	db.tableUsers = cfg.Config["db.table.users"]

	err = connectToDatabase()
	if err == nil {
		err = checkUsersTable()
	}

	return
}
