package db

import (
	loger "github.com/dredfort42/tools/logprinter"
)

// AddNewUser adds a new user to the database
func AddNewUser(email string, password string, accessToken string, refreshToken string) {
	query := `INSERT INTO ` + db.tableUsers + ` (
		email, 
		password_hash, 
		email_verified,
		remember_me,
		device_uuid,
		device_access_token, 
		access_token, 
		refresh_token, 
		created_at, 
		updated_at
	) VALUES (
		$1,
		crypt($2, gen_salt('bf')),
		FALSE, 
		FALSE,
		NULL,
		NULL,
		$3,
		$4,
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP
		)`

	if _, db.err = db.database.Exec(query, email, password, accessToken, refreshToken); db.err != nil {
		loger.Error("Failed to add new user to the database", db.err)
	} else {
		loger.Debug("New user successfully aded to the database")
	}
}
