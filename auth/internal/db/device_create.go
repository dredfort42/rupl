package db

import loger "github.com/dredfort42/tools/logprinter"

// DeviceCreate creates a new device
func DeviceCreate(email string, deviceUUID string, deviceAccessToken string, deviceRefreshToken string) (err error) {
	query := `
		INSERT INTO ` + db.tableDevices + `
			(device_uuid, device_access_token, device_refresh_token, email, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (device_uuid, email) DO 
			UPDATE SET
				device_access_token = $2,
				device_refresh_token = $3,
				created_at = CURRENT_TIMESTAMP;
	`

	_, err = db.database.Exec(query, deviceUUID, deviceAccessToken, deviceRefreshToken, email)
	if err != nil {
		loger.Error("Failed to create a new device", err)
	}

	return
}
