package db

import (
	"github.com/dredfort42/tools/logprinter"
)

// AddNewUser adds a new user to the database
func AddNewUser(email string, password string) {
	query := `INSERT INTO ` + db.tableUsers + ` (
		username, 
		password_hash, 
		email, 
		email_verified, 
		device_uuid, 
		access_token, 
		refresh_token, 
		created_at, 
		updated_at
	) VALUES (
		NULL,
		crypt($1, gen_salt('bf')), 
		$2,
		FALSE,
		NULL,
		NULL,
		NULL,
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP
		)`

	// username VARCHAR(100) NOT NULL,
	// 				password_hash VARCHAR(255) NOT NULL,
	// 				email VARCHAR(255) NOT NULL,
	// 				email_verified BOOLEAN DEFAULT FALSE,
	// 				device_uuid UUID NOT NULL,
	// 				access_token VARCHAR(255),
	// 				refresh_token VARCHAR(255),
	// 				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	// 				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

	_, db.err = db.database.Exec(query, password, email)

	if db.err != nil {
		logprinter.PrintError("Failed to add new user to the database", db.err)
	} else {
		logprinter.PrintSuccess("New user successfully aded to the database", "")
	}
}
