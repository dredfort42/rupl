package db

// AddNewUser adds a new user to the database
func AddNewUser(email string, password string, accessToken string, refreshToken string) (err error) {
	query := `
		INSERT INTO ` + db.tableUsers + ` (
			email, 
			password_hash, 
			email_verified,
			remember_me,
			access_token, 
			refresh_token, 
			devices,
			created_at, 
			updated_at
		) VALUES (
			$1,
			crypt($2, gen_salt('bf')),
			FALSE, 
			FALSE,
			$3,
			$4,
			NULL,
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP
		)`

	_, err = db.database.Exec(query, email, password, accessToken, refreshToken)

	return
}
