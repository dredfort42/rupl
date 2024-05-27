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
	for {
		_, err = db.database.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")
		if err != nil {
			return
		} else {
			loger.Success("Extension successfully created", "pgcrypto")
			break
		}
	}

	query := `
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_type 
				WHERE typname = 'devices'
			) THEN 
				CREATE TYPE devices AS (
					device_uuid UUID,
					device_access_token VARCHAR(255),
				);
			END IF;
		END $$
		CREATE TABLE IF NOT EXISTS ` + db.tableUsers + ` (
			email VARCHAR(255) PRIMARY KEY,
			password_hash VARCHAR(255) NOT NULL,
			remember_me BOOLEAN DEFAULT FALSE,
			email_verified BOOLEAN DEFAULT FALSE,
			access_token VARCHAR(255),
			refresh_token VARCHAR(255),
			devices devices[] DEFAULT '{}'::stride_length[] NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	if !isTableExists(db.tableUsers) {
		for {
			_, err = db.database.Exec(query)
			if err != nil {
				return
			} else {
				loger.Success("Table successfully created", db.tableUsers)
				break
			}
		}
	}

	return
}

// -- Insert user with hashed password
// INSERT INTO user_credentials (username, password_hash, device_uuid)
// VALUES ('example_user', crypt('user_password', gen_salt('bf')), 'device_uuid');

// -- Select user with hashed password
// SELECT * FROM user_credentials
// WHERE username = 'example_user'
// AND password_hash = crypt('user_input_password', password_hash);
