package db

import (
	"time"

	"github.com/dredfort42/tools/logprinter"
)

// CheckUsersTable checks if the users table exists, if not, it creates it
func CheckUsersTable() {
	tabalExists := false

	for !tabalExists {
		query := "SELECT * FROM " + db.tableUsers + ";"

		if _, db.err = db.database.Query(query); db.err != nil {
			extensionExists := false

			for !extensionExists {
				query = "CREATE EXTENSION IF NOT EXISTS pgcrypto;"
				if _, db.err = db.database.Exec(query); db.err != nil {
					logprinter.PrintError("Failed to create extension", db.err)
				} else {
					extensionExists = true
					logprinter.PrintSuccess("Extension successfully created", "pgcrypto")
				}
				time.Sleep(5 * time.Second)
			}

			logprinter.PrintWarning("Table does not exist", db.tableUsers)
			query = `
				CREATE TABLE IF NOT EXISTS ` + db.tableUsers + ` (
					id SERIAL PRIMARY KEY,
					username VARCHAR(100) NOT NULL,
					password_hash VARCHAR(255) NOT NULL,
					device_uuid UUID NOT NULL,
					access_token VARCHAR(255),
					refresh_token VARCHAR(255),
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
			`
			if _, db.err = db.database.Exec(query); db.err != nil {
				logprinter.PrintError("Failed to create table", db.err)
			} else {
				tabalExists = true
				logprinter.PrintSuccess("Table successfully created", db.tableUsers)
			}
		} else {
			tabalExists = true
			logprinter.PrintSuccess("Table found successfully", db.tableUsers)
		}
		time.Sleep(5 * time.Second)
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
