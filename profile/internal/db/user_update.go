package db

import (
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// UserUpdate updates a user profile in the database
func UserUpdate(user s.User) (err error) {
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
