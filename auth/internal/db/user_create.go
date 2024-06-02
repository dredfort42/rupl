package db

// UserSignUp adds a new user to the database
func UserSignUp(email string, password string, accessToken string, refreshToken string) (err error) {
	query := `
		INSERT INTO ` + db.tableUsers + ` (
			id,
			email,
			password_hash,
			access_token,
			refresh_token,
			user_browsers,
			user_devices,
			email_created_at,
			is_email_confirmed,
			created_at,
			updated_at
		) VALUES (
			DEFAULT,
			$1,
			crypt($2, gen_salt('bf')),
			$3,
			$4,
			DEFAULT,
			DEFAULT,
			CURRENT_TIMESTAMP,
			DEFAULT,
			DEFAULT,
			DEFAULT
		)`

	_, err = db.database.Exec(query, email, password, accessToken, refreshToken)

	return
}
