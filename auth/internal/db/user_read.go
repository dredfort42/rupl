package db

// IsUserExists checks if a user exists in the database
func IsUserExists(email string) (result bool) {
	query := `
		SELECT 1 
		FROM ` + db.tableUsers + ` 
		WHERE email = $1
	`

	db.database.QueryRow(query, email).Scan(&result)

	return
}

// IsUserPasswordCorrect checks if a user's password is correct
func IsUserPasswordCorrect(email string, password string) (result bool) {
	query := `
		SELECT 1
		FROM ` + db.tableUsers + `
		WHERE email = $1
		AND password_hash = crypt($2, password_hash)
	`

	db.database.QueryRow(query, email, password).Scan(&result)

	return
}
