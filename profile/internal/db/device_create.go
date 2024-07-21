package db

import (
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// DeviceCreate creates a new device in the database
func DeviceCreate(email string, device s.Device) (err error) {
	if DeviceExistsCheck(email, device.DeviceID) {
		DeviceDelete(device.DeviceID)
	}

	query := `
		INSERT INTO ` + db.tableDevices + ` (email, device_model, device_name, system_name, system_version, device_id, app_version) 
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	_, err = db.database.Exec(query, email, device.DeviceModel, device.DeviceName, device.SystemName, device.SystemVersion, device.DeviceID, device.AppVersion)
	if err != nil {
		loger.Error("Failed to create device in the database", err)
	}

	return
}
