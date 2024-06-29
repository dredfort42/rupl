package db

import loger "github.com/dredfort42/tools/logprinter"

// UserPasswordChange updates a user's password in the database
func UserPasswordChange(email string, newPassword string) (err error) {
	query := `
		UPDATE ` + db.tableUsers + `
		SET	password_hash = crypt($2, gen_salt('bf')),
			updated_at = CURRENT_TIMESTAMP
		WHERE email = $1;
	`

	_, err = db.database.Exec(query, email, newPassword)
	if err != nil {
		loger.Error("Failed to update user in the database", err)
	}

	return
}

// UserEmailChange updates a user's email address in the database
func UserEmailChange(email string, newEmail string) (err error) {
	query := `
		UPDATE ` + db.tableUsers + `
		SET	email = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE email = $1;
	`

	_, err = db.database.Exec(query, email, newEmail)
	if err != nil {
		loger.Error("Failed to update user in the database", err)
	}

	query = `
		UPDATE ` + db.tableDevices + `
		SET	email = $2
		WHERE email = $1;
	`

	_, err = db.database.Exec(query, email, newEmail)
	if err != nil {
		loger.Error("Failed to update devices in the database", err)
	}

	return
}
