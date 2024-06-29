package db

import loger "github.com/dredfort42/tools/logprinter"

// DeviceRefresh updates the device's access and refresh tokens
func DeviceRefresh(deviceUUID string, deviceAccessToken string, deviceRefreshToken string) (err error) {
	query := `
		UPDATE ` + db.tableDevices + `
		SET device_access_token = $2, device_refresh_token = $3, created_at = CURRENT_TIMESTAMP
		WHERE device_uuid = $1;
	`

	_, err = db.database.Exec(query, deviceUUID, deviceAccessToken, deviceRefreshToken)
	if err != nil {
		loger.Error("Failed to update device in the database", err)
	}

	return
}
