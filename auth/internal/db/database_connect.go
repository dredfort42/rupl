package db

import (
	"database/sql"
	"fmt"

	cfg "github.com/dredfort42/tools/configreader"
)

// databaseConnect connects to the database
func databaseConnect() {
	var url string = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Config["db.host"],
		cfg.Config["db.port"],
		cfg.Config["db.user"],
		cfg.Config["db.password"],
		cfg.Config["db.database.name"],
		cfg.Config["db.security.ssl"])

	var err error
	db.database, err = sql.Open("postgres", url)
	if err != nil {
		db.database.Close()
		panic(err)
	}

	err = db.database.Ping()
	if err != nil {
		db.database.Close()
		panic(err)
	}
}
