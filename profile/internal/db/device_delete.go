package db

import (
	loger "github.com/dredfort42/tools/logprinter"
)

// DeviceDelete deletes a device from the database
func DeviceDelete(deviceID string) (err error) {
	query := `
		DELETE FROM ` + db.tableDevices + ` 
		WHERE device_id = $1
	`

	_, err = db.database.Exec(query, deviceID)
	if err != nil {
		loger.Error("Failed to delete device from the database", err)
	}

	return
}
