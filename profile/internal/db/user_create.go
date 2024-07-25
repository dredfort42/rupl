package db

import (
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// UserCreate creates a new profile in the database
func UserCreate(user s.User) (err error) {
	query := `
		INSERT INTO ` + db.tableUsers + ` (email, first_name, last_name, date_of_birth, gender, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err = db.database.Exec(query, user.Email, user.FirstName, user.LastName, user.DateOfBirth, user.Gender)
	if err != nil {
		loger.Error("Failed to create profile in the database", err)
	}

	return
}
