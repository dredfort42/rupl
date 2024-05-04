package db

import (
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// CreateDevice creates a new device in the database
func CreateDevice(email string, device s.Device) error {
	if CheckDeviceExists(device.DeviceID) {
		DeleteDevice(device.DeviceID)
	}

	query := `INSERT INTO ` + db.tableDevices + ` (email, device_model, device_name, system_name, system_version, device_id, app_version) VALUES ($1, $2, $3, $4, $5, $6, $7);`

	if _, err := db.database.Exec(query, email, device.DeviceModel, device.DeviceName, device.SystemName, device.SystemVersion, device.DeviceID, device.AppVersion); err != nil {
		loger.Error("Failed to create device in the database", err)

		return err
	}

	return nil
}
