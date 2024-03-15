package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/dredfort42/tools/configreader"
	"github.com/dredfort42/tools/logprinter"
)

var db Database
var config configreader.ConfigMap

// ConnectToDatabase connects to the database and returns a pointer to it
func ConnectToDatabase() bool {
	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config["db.host"],
		config["db.port"],
		config["db.user"],
		config["db.password"],
		config["db.database.name"],
		config["db.security.ssl"])

	// logprinter.PrintInfo("Connecting to database", url)

	db.database, db.err = sql.Open("postgres", url)

	if db.err != nil {
		db.database.Close()
		logprinter.PrintError("Failed to connect to database", db.err)
		return false
	} else if db.err = db.database.Ping(); db.err != nil {
		db.database.Close()
		logprinter.PrintError("Failed to connect to database", db.err)
		return false
	}

	logprinter.PrintSuccess("Successfully connected to database", "")
	return true
}

var DEBUG bool = false

// Start starts the web service
func Start(configMap configreader.ConfigMap) {
	if debug := os.Getenv("DEBUG"); debug != "" {
		DEBUG = true
	}

	config = configMap

	db.database, db.err = nil, nil
	db.tableProfiles = config["db.table.profiles"]

	for !ConnectToDatabase() {
		time.Sleep(5 * time.Second)
	}

	CheckUsersTable()

	if DEBUG {
		logprinter.PrintSuccess("Database", "Started")
	}
}
