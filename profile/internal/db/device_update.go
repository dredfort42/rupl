package db

import (
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// DeviceUpdate updates a device in the database
func DeviceUpdate(email string, device s.Device) (err error) {
	query := `
		UPDATE ` + db.tableDevices + ` 
		SET device_model = $1, device_name = $2, system_name = $3, system_version = $4, app_version = $5, updated_at = CURRENT_TIMESTAMP 
		WHERE device_uuid = $6 AND email = $7;`

	_, err = db.database.Exec(query, device.DeviceModel, device.DeviceName, device.SystemName, device.SystemVersion, device.AppVersion, device.DeviceUUID, email)
	if err != nil {
		loger.Error("Failed to update device in the database", err)
	}

	return
}
