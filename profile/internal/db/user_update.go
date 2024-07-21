package db

import (
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// UserUpdate updates a user profile in the database
func UserUpdate(user s.User) (err error) {
	if !UserExistsCheck(user.Email) {
		loger.Warning("User does not exist in the database")
		return
	}

	query := `
		UPDATE ` + db.tableUsers + ` 
		SET first_name = $2, last_name = $3, date_of_birth = $4, gender = $5, updated_at = CURRENT_TIMESTAMP
		WHERE email = $1
	`

	_, err = db.database.Exec(query, user.Email, user.FirstName, user.LastName, user.DateOfBirth, user.Gender)
	if err != nil {
		loger.Error("Failed to update profile in the database", err)
	}

	return
}

// UserEmailChange updates a user's email address in the database
func UserEmailChange(email string, newEmail string) (err error) {
	if !UserExistsCheck(email) {
		loger.Warning("User does not exist in the database")
		return
	}

	query := `
		UPDATE ` + db.tableUsers + ` 
		SET email = $2, 
			updated_at = CURRENT_TIMESTAMP
		WHERE email = $1;
	`

	_, err = db.database.Exec(query, email, newEmail)
	if err != nil {
		loger.Error("Failed to update email in the users table", err)
	}

	if !UserDevicesExistsCheck(email) {
		loger.Warning("Device does not exist in the database")
		return
	}

	query = `
		UPDATE ` + db.tableDevices + `
		SET email = $2, 
			updated_at = CURRENT_TIMESTAMP
		WHERE email = $1;
	`

	_, err = db.database.Exec(query, email, newEmail)
	if err != nil {
		loger.Error("Failed to update email in the devices table", err)
	}

	return
}
