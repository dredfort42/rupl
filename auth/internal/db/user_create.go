package db

// UserSignUp adds a new user to the database
func UserSignUp(email string, password string, accessToken string, refreshToken string) (err error) {
	query := `
		INSERT INTO ` + db.tableUsers + ` (
			email,
			password_hash
		) VALUES (
			$1,
			crypt($2, gen_salt('bf'))
		)
	`

	_, err = db.database.Exec(query, email, password)
	if err != nil {
		return
	}

	err = SessionCreate(email, accessToken, refreshToken, true)

	return
}
