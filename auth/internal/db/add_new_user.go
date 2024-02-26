package db

import (
	"github.com/dredfort42/tools/logprinter"
)

// AddNewUser adds a new user to the database
func AddNewUser(email string, password string) {
	query := `INSERT INTO ` + db.tableUsers + ` (
		username, 
		email, 
		password_hash, 
		email_verified, 
		device_uuid, 
		access_token, 
		refresh_token, 
		created_at, 
		updated_at
	) VALUES (
		NULL,
		$1,
		crypt($2, gen_salt('bf')), 
		FALSE,
		NULL,
		NULL,
		NULL,
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP
		)`

	if _, db.err = db.database.Exec(query, email, password); db.err != nil {
		logprinter.PrintError("Failed to add new user to the database", db.err)
	} else {
		logprinter.PrintSuccess("New user successfully aded to the database", "")
	}
}
