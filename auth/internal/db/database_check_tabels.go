package db

import (
	loger "github.com/dredfort42/tools/logprinter"
)

// isTableExists checks if the table exists
func isTableExists(tabelName string) (isTableExists bool) {
	err := db.database.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", tabelName).Scan(&isTableExists)
	if err != nil || !isTableExists {
		loger.Warning("Table does not exist", tabelName)
	} else {
		loger.Debug("Table found successfully", tabelName)
	}

	return
}

// checkUsersTable checks if the users table exists, if not, it creates it
func checkUsersTable() (err error) {
	if isTableExists(db.tableUsers) {
		return
	}

	_, err = db.database.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	if err != nil {
		loger.Debug(err.Error())
		return
	}

	query := `
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_type 
				WHERE typname = 'user_browsers'
			) THEN
				CREATE TYPE user_browsers AS (
					remembered_access_token VARCHAR(255) NOT NULL,
					remembered_refresh_token VARCHAR(255) NOT NULL
				);
			END IF;
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_type 
				WHERE typname = 'user_devices'
			) THEN
				CREATE TYPE user_devices AS (
					device_uuid UUID NOT NULL,
					device_access_token VARCHAR(255) NOT NULL
				);
			END IF;
		END $$;

		CREATE TABLE IF NOT EXISTS ` + db.tableUsers + ` (
			email VARCHAR(255) PRIMARY KEY,
			password_hash VARCHAR(255) NOT NULL,
			access_token VARCHAR(255) NOT NULL,
			refresh_token VARCHAR(255) NOT NULL,
			user_browsers user_browsers[] DEFAULT '{}'::user_browsers[] NOT NULL,
			user_devices user_devices[] DEFAULT '{}'::user_devices[] NOT NULL,
			email_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			is_email_confirmed BOOLEAN DEFAULT FALSE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
		);
	`

	_, err = db.database.Exec(query)
	if err != nil {
		loger.Debug(err.Error())
		return
	} else {
		loger.Success("Table successfully created", db.tableUsers)
	}

	return
}
