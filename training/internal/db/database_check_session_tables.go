package db

import (
	"time"

	loger "github.com/dredfort42/tools/logprinter"
)

// checkSessionTables() checks if the sessions table exists, if not, it creates it
func checkSessionsTotalTable() {
	var tabalExists bool = checkTableExists(db.tableSessions + "_total")

	for !tabalExists {
		query := `
				CREATE TABLE IF NOT EXISTS ` + db.tableSessions + `_total (
					session_uuid VARCHAR(255) PRIMARY KEY,
					email VARCHAR(255) NOT NULL,
					timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
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

// checkSessionsRouteTable() checks if the sessions table exists, if not, it creates it
func checkSessionsRouteTable() {
	var tabalExists bool = checkTableExists(db.tableSessions + "_route")

	for !tabalExists {
		query := `
			DO $$ BEGIN
				IF NOT EXISTS (
					SELECT 1 
					FROM pg_type 
					WHERE typname = 'route_data'
				) THEN 
					CREATE TYPE route_data AS (
						timestamp BIGINT,
						latitude DOUBLE PRECISION,
						longitude DOUBLE PRECISION,
						horizontal_accuracy DOUBLE PRECISION,
						altitude DOUBLE PRECISION,
						vertical_accuracy DOUBLE PRECISION,
						speed DOUBLE PRECISION,
						speed_accuracy DOUBLE PRECISION,
						course DOUBLE PRECISION,
						course_accuracy DOUBLE PRECISION
					);
				END IF;
			END $$;
			CREATE TABLE IF NOT EXISTS public.` + db.tableSessions + `_route (
				session_uuid VARCHAR(255) NOT NULL,
				route_data route_data[] DEFAULT '{}'::route_data[] NOT NULL,
				constraint ` + db.tableSessions + `_route_pk PRIMARY KEY (session_uuid)
			);
		`

		if _, db.err = db.database.Exec(query); db.err != nil {
			loger.Error("Failed to create table", db.err)
			time.Sleep(5 * time.Second)
		} else {
			tabalExists = true
			loger.Success("Table successfully created", db.tableSessions+"_route")
		}
	}
}

func checkSessionTables() {
	checkSessionsTotalTable()
	checkSessionsRouteTable()
}
