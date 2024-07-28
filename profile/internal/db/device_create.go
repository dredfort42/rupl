package db

import (
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// DeviceCreate creates a new device in the database
func DeviceCreate(email string, device s.Device) (err error) {
	if DeviceExistsCheck(email, device.DeviceUUID) {
		DeviceDelete(email, device.DeviceUUID)
	}

	query := `
		INSERT INTO ` + db.tableDevices + ` (email, device_model, device_name, system_name, system_version, device_uuid, app_version) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (email, device_uuid) DO UPDATE SET
			device_model = $2,
			device_name = $3,
			system_name = $4,
			system_version = $5,
			app_version = $7,
			updated_at = CURRENT_TIMESTAMP;
	`

	_, err = db.database.Exec(query, email, device.DeviceModel, device.DeviceName, device.SystemName, device.SystemVersion, device.DeviceUUID, device.AppVersion)
	if err != nil {
		loger.Error("Failed to create device in the database", err)
	}

	return
}
