package db

import loger "github.com/dredfort42/tools/logprinter"

// SessionUpdateOneTime updates token for one-time access
func SessionUpdateOneTime(email string, accessToken string, refreshToken string) (err error) {
	err = SessionDeleteOneTime(email)
	if err != nil {
		return
	}

	err = SessionCreate(email, accessToken, refreshToken, true)
	if err != nil {
		loger.Error("Failed to update one-time session in the database", err)
	}

	return
}

// SessionUpdate updates access and refresh tokens
func SessionUpdate(email string, refreshToken string, newAccessToken string, newRefreshToken string) (err error) {
	query := `
		UPDATE ` + db.tableSessions + `
		SET access_token = $3, refresh_token = $4, created_at = CURRENT_TIMESTAMP
		WHERE email = $1
		AND refresh_token = $2;
	`

	_, err = db.database.Exec(query, email, refreshToken, newAccessToken, newRefreshToken)
	if err != nil {
		loger.Error("Failed to update session in the database", err)
	}

	return
}
