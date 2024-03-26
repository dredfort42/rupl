package db

import (
	"github.com/dredfort42/tools/logprinter"
)

// Delete user's access and refresh tokens from the database
func DeleteUserTokens(email string) {
	query := `UPDATE ` + db.tableUsers + ` SET access_token = NULL, refresh_token = NULL WHERE email = $1`

	if _, db.err = db.database.Exec(query, email); db.err != nil {
		logprinter.PrintError("Failed to delete "+email+" access and refresh tokens from the database", db.err)
	} else {
		logprinter.PrintSuccess(email+" user's access and refresh tokens successfully deleted from the database", "")
	}
}

// Delete device access token from the database
func DeleteDeviceAccessToken(clientID string, accessToken string) bool {
	query := `UPDATE FROM ` + db.tableUsers + ` SET device_uuid = NULL, device_access_token = NULL WHERE device_uuid = $1 AND device_access_token = $2`

	if _, db.err = db.database.Exec(query, clientID, accessToken); db.err != nil {
		logprinter.PrintError("Failed to delete device access token from the database", db.err)
		return false
	}

	logprinter.PrintSuccess("Device access token successfully deleted from the database", "")
	return true
}
