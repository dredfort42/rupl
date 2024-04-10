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

// CheckUsersTable checks if the users table exists, if not, it creates it
func CheckUsersTable() {
	var tabalExists bool = CheckTableExists(db.tableUsers)

	for !tabalExists {

		extensionExists := false

		for !extensionExists {
			query := "CREATE EXTENSION IF NOT EXISTS pgcrypto;"
			if _, db.err = db.database.Exec(query); db.err != nil {
				logprinter.PrintError("Failed to create extension", db.err)
			} else {
				extensionExists = true
				logprinter.PrintSuccess("Extension successfully created", "pgcrypto")
			}
			time.Sleep(5 * time.Second)
		}

		query := `
				CREATE TABLE IF NOT EXISTS ` + db.tableUsers + ` (
					email VARCHAR(255) PRIMARY KEY,
					password_hash VARCHAR(255) NOT NULL,
					remember_me BOOLEAN DEFAULT FALSE,
					email_verified BOOLEAN DEFAULT FALSE,
					device_uuid UUID,
					device_access_token VARCHAR(255),
					access_token VARCHAR(255),
					refresh_token VARCHAR(255),
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
			`
		if _, db.err = db.database.Exec(query); db.err != nil {
			logprinter.PrintError("Failed to create table", db.err)
			time.Sleep(5 * time.Second)
		} else {
			tabalExists = true
			logprinter.PrintSuccess("Table successfully created", db.tableUsers)
		}
	}
}

// CheckTables checks if the tables exists, if not, it creates it
func CheckTables() {
	CheckUsersTable()
}

// -- Insert user with hashed password
// INSERT INTO user_credentials (username, password_hash, device_uuid)
// VALUES ('example_user', crypt('user_password', gen_salt('bf')), 'device_uuid');

// -- Select user with hashed password
// SELECT * FROM user_credentials
// WHERE username = 'example_user'
// AND password_hash = crypt('user_input_password', password_hash);
