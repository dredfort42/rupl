package db

import (
	loger "github.com/dredfort42/tools/logprinter"
)

// DeleteDevice deletes a device from the database
func DeleteDevice(deviceID string) error {
	query := `DELETE FROM ` + db.tableDevices + ` WHERE device_id = $1`

	if _, err := db.database.Exec(query, deviceID); err != nil {
		loger.Error("Failed to delete device from the database", err)

		return err
	}

	return nil
}
