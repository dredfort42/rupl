package db

import (
	"errors"

	"github.com/dredfort42/tools/logprinter"
)

// CreateProfile creates a new profile in the database
func CreateProfile(profile Profile) error {
	if CheckProfileExists(profile.Email) {
		return errors.New("Profile already exists")
	}

	query := `INSERT INTO ` + db.tableProfiles + ` (email, first_name, last_name, date_of_birth, gender, created_at, 
		updated_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	if _, err := db.database.Exec(query, profile.Email, profile.FirstName, profile.LastName, profile.DateOfBirth, profile.Gender); err != nil {
		if DEBUG {
			logprinter.PrintError("Failed to create profile in the database", err)
		}
		return err
	}
	return nil
}
