package db

import (
	"github.com/dredfort42/tools/logprinter"
)

// UpdateDevice updates a device in the database
func UpdateDevice(device Device) error {
	query := `UPDATE ` + db.tableDevices + ` SET device_model = $1, device_name = $2, system_name = $3, system_version = $4, app_version = $5 WHERE device_id = $6;`
	if _, err := db.database.Exec(query, device.DeviceModel, device.DeviceName, device.SystemName, device.SystemVersion, device.AppVersion, device.DeviceID); err != nil {
		if DEBUG {
			logprinter.PrintError("Failed to update device in the database", err)
		}
		return err
	}
	return nil
}
