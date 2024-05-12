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

// checkSessionTables() checks if the sessions table exists, if not, it creates it
func checkSessionsTotalTable() {
	var tabalExists bool = checkTableExists(db.tableSessions + "_total")

	for !tabalExists {
		query := `
				CREATE TABLE IF NOT EXISTS ` + db.tableSessions + `_total (
					session_uuid VARCHAR(255) PRIMARY KEY,
					email VARCHAR(255) NOT NULL,
					timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
					steps INT NOT NULL DEFAULT 0,
					vo2max FLOAT NOT NULL DEFAULT 0,
					avr_power FLOAT NOT NULL DEFAULT 0,
					avr_oscillation FLOAT NOT NULL DEFAULT 0,
					avr_contact_time FLOAT NOT NULL DEFAULT 0,
					avr_heartrate FLOAT NOT NULL DEFAULT 0,
					avr_stride_length FLOAT NOT NULL DEFAULT 0,
					avr_ground_contact_time FLOAT NOT NULL DEFAULT 0,
					avr_speed FLOAT NOT NULL DEFAULT 0,
					total_flights_climbed INT NOT NULL DEFAULT 0,
					total_energy_burned INT NOT NULL DEFAULT 0,
					total_distance FLOAT NOT NULL DEFAULT 0,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
			`
		if _, db.err = db.database.Exec(query); db.err != nil {
			loger.Error("Failed to create table", db.err)
			time.Sleep(5 * time.Second)
		} else {
			tabalExists = true
			loger.Success("Table successfully created", db.tableSessions+"_total")
		}
	}
}

func checkSessionTables() {
	checkSessionsTotalTable()
}

// CheckTables checks if the tables exists, if not, it creates it
func checkTables() {
	checkSessionTables()
}
