package db

import (
	"errors"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// ProfileCreate creates a new profile in the database
func ProfileCreate(profile s.Profile) (err error) {
	if CheckProfileExists(profile.Email) {
		return errors.New("profile already exists")
	}

	query := `
		INSERT INTO ` + db.tableUsers + ` (email, first_name, last_name, date_of_birth, gender, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err = db.database.Exec(query, profile.Email, profile.FirstName, profile.LastName, profile.DateOfBirth, profile.Gender)
	if err != nil {
		loger.Error("Failed to create profile in the database", err)
	}

	return
}
