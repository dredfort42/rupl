package db

import (
	loger "github.com/dredfort42/tools/logprinter"
)

// Delete user's access and refresh tokens from the database
func DeleteUserTokens(email string) bool {
	query := `UPDATE ` + db.tableUsers + ` SET access_token = NULL, refresh_token = NULL WHERE email = $1`

	if _, db.err = db.database.Exec(query, email); db.err != nil {
		loger.Error("Failed to delete "+email+" access and refresh tokens from the database", db.err)
		return false
	}

	loger.Debug(email+" user's access and refresh tokens successfully deleted from the database", "")
	return true
}

// Delete device access token from the database
func DeleteDeviceAccessToken(email string) bool {
	query := `UPDATE ` + db.tableUsers + ` SET device_uuid = NULL, device_access_token = NULL WHERE email = $1`

	if _, db.err = db.database.Exec(query, email); db.err != nil {
		loger.Error("Failed to delete "+email+" device access token from the database", db.err)
		return false
	}

	loger.Debug(email+" user's device access token successfully deleted from the database", "")
	return true
}
