package db

import (
	"database/sql"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// UserExistsCheck checks if a user exists in the database based on the email provided
func UserExistsCheck(email string) (result bool) {
	query := `
		SELECT 1
	 	FROM ` + db.tableUsers + ` 
		WHERE email = $1
	`

	err := db.database.QueryRow(query, email).Scan(&result)
	if err != nil && err != sql.ErrNoRows {
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
		if err == sql.ErrNoRows {
			err = nil
		} else {
			loger.Error("Failed to get profile from the database", err)
		}
	}

	return
}
