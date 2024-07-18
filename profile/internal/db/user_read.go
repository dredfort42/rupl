package db

import (
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// CheckUserExists checks if a user exists in the database based on the email provided
func CheckUserExists(email string) (result bool) {
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

// UserGet returns a user from the database
func UserGet(email string) (user s.User, err error) {
	query := `
		SELECT email, first_name, last_name, date_of_birth, gender
		FROM ` + db.tableUsers + ` 
		WHERE email = $1
	`

	err = db.database.QueryRow(query, email).Scan(&user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Gender)
	if err != nil {
		loger.Error("Failed to get profile from the database", err)
	}

	return
}
