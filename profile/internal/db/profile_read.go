package db

import (
	"github.com/dredfort42/tools/logprinter"
)

// CheckProfileExists checks if a profile exists in the database based on the email provided
func CheckProfileExists(email string) bool {
	query := `SELECT email FROM ` + db.tableProfiles + ` WHERE email = $1`

	if err := db.database.QueryRow(query, email).Scan(&email); err != nil {
		if DEBUG {
			logprinter.PrintError("Failed to check if profile exists in the database", err)
		}
		return false
	}
	return true
}

// GetProfile returns a profile from the database
func GetProfile(email string) (Profile, error) {
	query := `SELECT * FROM ` + db.tableProfiles + ` WHERE email = $1`

	var profile Profile
	var created_at string
	var updated_at string

	if err := db.database.QueryRow(query, email).Scan(&profile.Email, &profile.FirstName, &profile.LastName, &profile.DateOfBirth, &profile.Gender, &created_at, &updated_at); err != nil {
		if DEBUG {
			logprinter.PrintError("Failed to get profile from the database", err)
		}
		return profile, err
	}
	return profile, nil
}
