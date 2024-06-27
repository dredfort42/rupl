package db

import (
	loger "github.com/dredfort42/tools/logprinter"
)

// SessionCreate adds a new session to the database
func SessionCreate(email string, accessToken string, refreshToken string, isOneTime bool) (err error) {
	err = SessionDeleteOneTime(email)
	if err != nil {
		return
	}

	query := `
		INSERT INTO ` + db.tableSessions + ` (
			email,
			access_token,
			refresh_token,
			is_one_time,
			created_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			CURRENT_TIMESTAMP
		)
	`

	_, err = db.database.Exec(query, email, accessToken, refreshToken, isOneTime)
	if err != nil {
		loger.Error("Failed to create session in the database", err)
	}

	return
}
