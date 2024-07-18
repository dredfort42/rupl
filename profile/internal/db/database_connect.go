package db

import (
	"database/sql"

	cfg "github.com/dredfort42/tools/configreader"
)

// Database is the database struct
type Database struct {
	database     *sql.DB
	tableUsers   string
	tableDevices string
}

var db Database

// databaseConnect connects to the database and returns a pointer to it
func databaseConnect() {
	host := cfg.Config["db.host"]
	if host == "" {
		panic("Database host is not set")
	}

	port := cfg.Config["db.port"]
	if port == "" {
		panic("Database port is not set")
	}

	user := cfg.Config["db.user"]
	if user == "" {
		panic("Database user is not set")
	}

	password := cfg.Config["db.password"]
	if password == "" {
		panic("Database password is not set")
	}

	databaseName := cfg.Config["db.database.name"]
	if databaseName == "" {
		panic("Database name is not set")
	}

	ssl := cfg.Config["db.security.ssl"]
	if ssl == "" {
		panic("Database ssl is not set")
	}

	url := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + databaseName + " sslmode=" + ssl

	var err error
	db.database, err = sql.Open("postgres", url)
	if err != nil {
		db.database.Close()
		panic("Database connection error | " + err.Error())
	}

	err = db.database.Ping()
	if err != nil {
		db.database.Close()
		panic("Database ping error | " + err.Error())
	}
}
