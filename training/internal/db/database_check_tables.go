package db

import (
	loger "github.com/dredfort42/tools/logprinter"
)

// checkTableExists checks if the table exists
func checkTableExists(tabelName string) bool {
	tabelExists := false

	err := db.database.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", tabelName).Scan(&tabelExists)
	if err != nil || !tabelExists {
		loger.Warning("Table does not exist", tabelName)
		return false
	} else {
		loger.Debug("Table found successfully", tabelName)
		return true
	}
}

// CheckTables checks if the tables exists, if not, it creates it
func checkTables() {
	checkSessionTables()
}
