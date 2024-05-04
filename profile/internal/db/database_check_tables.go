package db

import (
	"time"

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

// checkProfilesTable() checks if the users table exists, if not, it creates it
func checkProfilesTable() {
	var tabalExists bool = checkTableExists(db.tableProfiles)

	for !tabalExists {
		query := `
				CREATE TABLE IF NOT EXISTS ` + db.tableProfiles + ` (
					email VARCHAR(255) PRIMARY KEY,
					first_name VARCHAR(255) NOT NULL,
					last_name VARCHAR(255) NOT NULL,
					date_of_birth DATE NOT NULL,
					gender VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
			`
		if _, db.err = db.database.Exec(query); db.err != nil {
			loger.Error("Failed to create table", db.err)
			time.Sleep(5 * time.Second)
		} else {
			tabalExists = true
			loger.Success("Table successfully created", db.tableProfiles)
		}
	}
}

// checkDevicesTable checks if the devices table exists, if not, it creates it
func checkDevicesTable() {
	var tabalExists bool = checkTableExists(db.tableDevices)

	for !tabalExists {
		query := `
				CREATE TABLE IF NOT EXISTS ` + db.tableDevices + ` (
					id SERIAL PRIMARY KEY,
					email VARCHAR(255) NOT NULL,
					device_model VARCHAR(255) NOT NULL,
					device_name VARCHAR(255) NOT NULL,
					system_name VARCHAR(255) NOT NULL,
					system_version VARCHAR(255) NOT NULL,
					device_id VARCHAR(255) NOT NULL,
					app_version VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
			`
		if _, db.err = db.database.Exec(query); db.err != nil {
			loger.Error("Failed to create table", db.err)
			time.Sleep(5 * time.Second)
		} else {
			tabalExists = true
			loger.Success("Table successfully created", db.tableDevices)
		}
	}
}

// CheckTables checks if the tables exists, if not, it creates it
func checkTables() {
	checkProfilesTable()
	checkDevicesTable()
}
