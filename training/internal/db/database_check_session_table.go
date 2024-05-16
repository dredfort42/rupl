package db

import (
	"time"

	loger "github.com/dredfort42/tools/logprinter"
)

// checkSessionTable() checks if the sessions table exists, if not, it creates it
func checkSessionTable() {
	var tabalExists bool = checkTableExists(db.tableSessions + "_total")

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
					IF NOT EXISTS (
						SELECT 1 
						FROM pg_type 
						WHERE typname = 'step_count'
					) THEN 
						CREATE TYPE step_count AS (
							timestamp BIGINT,
							steps_count INT
						);
					END IF;
					IF NOT EXISTS (
						SELECT 1 
						FROM pg_type 
						WHERE typname = 'running_power'
					) THEN 
						CREATE TYPE running_power AS (
							timestamp BIGINT,
							power_W INT
						);
					END IF;
					IF NOT EXISTS (
						SELECT 1 
						FROM pg_type 
						WHERE typname = 'vertical_oscillation'
					) THEN 
						CREATE TYPE vertical_oscillation AS (
							timestamp BIGINT,
							vertical_oscillation_cm FLOAT
						);
					END IF;
					IF NOT EXISTS (
						SELECT 1 
						FROM pg_type 
						WHERE typname = 'energy_burned'
					) THEN 
						CREATE TYPE energy_burned AS (
							timestamp BIGINT,
							energy_burned_kcal FLOAT
						);
					END IF;
					IF NOT EXISTS (
						SELECT 1 
						FROM pg_type 
						WHERE typname = 'heart_rate'
					) THEN 
						CREATE TYPE heart_rate AS (
							timestamp BIGINT,
							heart_rate_count_per_s FLOAT
						);
					END IF;
					IF NOT EXISTS (
						SELECT 1 
						FROM pg_type 
						WHERE typname = 'stride_length'
					) THEN 
						CREATE TYPE stride_length AS (
							timestamp BIGINT,
							stride_length_m FLOAT
						);
					END IF;
					IF NOT EXISTS (
						SELECT 1 
						FROM pg_type 
						WHERE typname = 'ground_contact_time'
					) THEN 
						CREATE TYPE ground_contact_time AS (
							timestamp BIGINT,
							ground_contact_time_ms INT
						);
					END IF;
					IF NOT EXISTS (
						SELECT 1 
						FROM pg_type 
						WHERE typname = 'speed'
					) THEN 
						CREATE TYPE speed AS (
							timestamp BIGINT,
							speed_m_per_s FLOAT
						);
					END IF;
					IF NOT EXISTS (
						SELECT 1 
						FROM pg_type 
						WHERE typname = 'distance'
					) THEN 
						CREATE TYPE distance AS (
							timestamp BIGINT,
							distance_m FLOAT
						);
					END IF;
				END $$;
				CREATE TABLE IF NOT EXISTS ` + db.tableSessions + ` (
					session_uuid VARCHAR(255) NOT NULL,
					session_start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					session_end_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					email VARCHAR(255) NOT NULL,
					route_data route_data[] DEFAULT '{}'::route_data[] NOT NULL,
					step_count step_count[] DEFAULT '{}'::step_count[] NOT NULL,
					running_power running_power[] DEFAULT '{}'::running_power[] NOT NULL,
					vertical_oscillation vertical_oscillation[] DEFAULT '{}'::vertical_oscillation[] NOT NULL,
					energy_burned energy_burned[] DEFAULT '{}'::energy_burned[] NOT NULL,
					heart_rate heart_rate[] DEFAULT '{}'::heart_rate[] NOT NULL,
					stride_length stride_length[] DEFAULT '{}'::stride_length[] NOT NULL,
					ground_contact_time ground_contact_time[] DEFAULT '{}'::ground_contact_time[] NOT NULL,
					speed speed[] DEFAULT '{}'::speed[] NOT NULL,
					distance distance[] DEFAULT '{}'::distance[] NOT NULL,
					vo2max_mL_per_min_per_kg FLOAT NOT NULL DEFAULT 0,
					avr_speed_m_per_s FLOAT NOT NULL DEFAULT 0,
					avr_heart_rate_count_per_s FLOAT NOT NULL DEFAULT 0,
					avr_power_W FLOAT NOT NULL DEFAULT 0,
					avr_vertical_oscillation_cm FLOAT NOT NULL DEFAULT 0,
					avr_ground_contact_time_ms FLOAT NOT NULL DEFAULT 0,
					avr_stride_length_m FLOAT NOT NULL DEFAULT 0,
					total_distance_m FLOAT NOT NULL DEFAULT 0,
					total_steps_count INT NOT NULL DEFAULT 0,
					total_energy_burned_kcal INT NOT NULL DEFAULT 0,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					constraint ` + db.tableSessions + `_pkey PRIMARY KEY (session_uuid)
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
