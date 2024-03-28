package db

import (
	"time"

	"github.com/dredfort42/tools/logprinter"
)

// CheckProfilesTable checks if the users table exists, if not, it creates it
func CheckProfilesTable() {
	tabalExists := false

	for !tabalExists {
		query := "SELECT * FROM " + db.tableProfiles + ";"

		if _, db.err = db.database.Query(query); db.err != nil {
			logprinter.PrintWarning("Table does not exist", db.tableProfiles)
			query = `
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
			} else {
				tabalExists = true
				logprinter.PrintSuccess("Table successfully created", db.tableProfiles)
			}
		} else {
			tabalExists = true
			logprinter.PrintSuccess("Table found successfully", db.tableProfiles)
		}
		time.Sleep(5 * time.Second)
	}
}

// CheckDevicesTable checks if the devices table exists, if not, it creates it
func CheckDevicesTable() {
	tabalExists := false

	for !tabalExists {
		query := "SELECT * FROM " + db.tableDevices + ";"

		if _, db.err = db.database.Query(query); db.err != nil {
			logprinter.PrintWarning("Table does not exist", db.tableDevices)
			query = `
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
			} else {
				tabalExists = true
				logprinter.PrintSuccess("Table successfully created", db.tableDevices)
			}
		} else {
			tabalExists = true
			logprinter.PrintSuccess("Table found successfully", db.tableDevices)
		}
		time.Sleep(5 * time.Second)
	}
}

// CheckTables checks if the tables exists, if not, it creates it
func CheckTables() {
	CheckProfilesTable()
	CheckDevicesTable()
}
