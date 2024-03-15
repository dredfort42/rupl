package db

import (
	"time"

	"github.com/dredfort42/tools/logprinter"
)

// CheckUsersTable checks if the users table exists, if not, it creates it
func CheckUsersTable() {
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
