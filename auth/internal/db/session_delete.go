package db

import (
	loger "github.com/dredfort42/tools/logprinter"
)

// SessionDeleteOneTime deletes one-time session from the database
func SessionDeleteOneTime(email string) (err error) {
	query := `
		DELETE FROM ` + db.tableSessions + `
		WHERE email = $1
		AND is_one_time = TRUE;
	`

	_, err = db.database.Exec(query, email)
	if err != nil {
		loger.Error("Failed to delete session from the database", err)
	}

	return
}

// SessionDelete deletes a session from the database
func SessionDelete(email string, accessToken string) (err error) {
	query := `
		DELETE FROM ` + db.tableSessions + `
		WHERE email = $1
		AND access_token = $2;
	`

	_, err = db.database.Exec(query, email, accessToken)
	if err != nil {
		loger.Error("Failed to delete session from the database", err)
	}

	return
}

// SessionDeleteAll deletes all sessions from the database
func SessionDeleteAll(email string) (err error) {
	query := `
		DELETE FROM ` + db.tableSessions + `
		WHERE email = $1;
	`

	_, err = db.database.Exec(query, email)
	if err != nil {
		loger.Error("Failed to delete all sessions from the database", err)
	}

	return
}
