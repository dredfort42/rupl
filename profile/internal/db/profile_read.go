package db

import (
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// CheckProfileExists checks if a profile exists in the database based on the email provided
func CheckProfileExists(email string) (result bool) {
	query := `
		SELECT 1
	 	FROM ` + db.tableUsers + ` 
		WHERE email = $1
	`

	err := db.database.QueryRow(query, email).Scan(&result)
	if err != nil {
		loger.Error("Failed to check if profile exists in the database", err)
	}

	return
}

// ProfileGet returns a profile from the database
func ProfileGet(email string) (profile s.Profile, err error) {
	query := `
		SELECT * 
		FROM ` + db.tableUsers + ` 
		WHERE email = $1
	`

	err = db.database.QueryRow(query, email).Scan(&profile.Email, &profile.FirstName, &profile.LastName, &profile.DateOfBirth, &profile.Gender)
	if err != nil {
		loger.Error("Failed to get profile from the database", err)
	}

	return
}
