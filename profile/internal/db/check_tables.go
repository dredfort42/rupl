package db

import (
	"time"

	"github.com/dredfort42/tools/logprinter"
)

// CheckTableExists checks if the table exists
func CheckTableExists(tabelName string) bool {
	tabelExists := false

	err := db.database.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", tabelName).Scan(&tabelExists)
	if err != nil || !tabelExists {
		logprinter.PrintWarning("Table does not exist", tabelName)
		return false
	} else {
		logprinter.PrintSuccess("Table found successfully", tabelName)
		return true
	}
}

// CheckProfilesTable checks if the users table exists, if not, it creates it
func CheckProfilesTable() {
	var tabalExists bool = CheckTableExists(db.tableProfiles)

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
			logprinter.PrintError("Failed to create table", db.err)
			time.Sleep(5 * time.Second)
		} else {
			tabalExists = true
			logprinter.PrintSuccess("Table successfully created", db.tableProfiles)
		}
	}
}

// CheckDevicesTable checks if the devices table exists, if not, it creates it
func CheckDevicesTable() {
	var tabalExists bool = CheckTableExists(db.tableDevices)

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
			logprinter.PrintError("Failed to create table", db.err)
			time.Sleep(5 * time.Second)
		} else {
			tabalExists = true
			logprinter.PrintSuccess("Table successfully created", db.tableDevices)
		}
	}
}

// CheckTables checks if the tables exists, if not, it creates it
func CheckTables() {
	CheckProfilesTable()
	CheckDevicesTable()
}
