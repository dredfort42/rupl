package db

import (
	// "errors"

	"github.com/dredfort42/tools/logprinter"
)

// CreateDevice creates a new device in the database
func CreateDevice(email string, device Device) error {
	// if CheckdeviceExists(device.Email) {
	// 	return errors.New("device already exists")
	// }

	query := `INSERT INTO ` + db.tableDevices + ` (email, device_model, device_name, system_name, system_version, device_id) VALUES ($1, $2, $3, $4, $5, $6);`
	if _, err := db.database.Exec(query, email, device.DeviceModel, device.DeviceName, device.SystemName, device.SystemVersion, device.DeviceID); err != nil {
		if DEBUG {
			logprinter.PrintError("Failed to create device in the database", err)
		}
		return err
	}
	return nil
}
