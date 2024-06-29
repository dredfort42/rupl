package db

import loger "github.com/dredfort42/tools/logprinter"

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
		loger.Error("Failed to create user in the database", err)
		return
	}

	err = SessionCreate(email, accessToken, refreshToken, true)
	if err != nil {
		loger.Error("Failed to create session in the database", err)
	}

	return
}
