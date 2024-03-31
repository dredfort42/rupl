package db

import (
	"github.com/dredfort42/tools/logprinter"
)

// DeleteDevice deletes a device from the database
func DeleteDevice(email string, deviceID string) error {
	query := `DELETE FROM ` + db.tableDevices + ` WHERE email = $1 AND device_id = $2`

	if _, err := db.database.Exec(query, email, deviceID); err != nil {
		if DEBUG {
			logprinter.PrintError("Failed to delete device from the database", err)
		}
		return err
	}
	return nil
}
