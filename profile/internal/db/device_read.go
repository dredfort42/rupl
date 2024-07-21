package db

import (
	"database/sql"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// CheckDeviceExists checks if a device exists in the database based on the email and device ID provided
func CheckDeviceExists(deviceID string) bool {
	query := `SELECT device_id FROM ` + db.tableDevices + ` WHERE device_id = $1;`

	err := db.database.QueryRow(query, deviceID).Scan(&deviceID)
	if err != nil {
		if err != sql.ErrNoRows {
			loger.Error("Failed to check if device exists in the database", err)
		}

		return false
	}

	return true
}

// DevicesGet returns a device from the database
func DevicesGet(email string) (devices s.UserDevices, err error) {
	query := `
		SELECT * 
		FROM ` + db.tableDevices + ` 
		WHERE email = $1;
	`

	rows, err := db.database.Query(query, email)
	if err != nil {
		return s.UserDevices{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var tmpEmail string
		var device s.Device
		var created_at string
		var updated_at string

		err = rows.Scan(&id, &tmpEmail, &device.DeviceModel, &device.DeviceName, &device.SystemName, &device.SystemVersion, &device.DeviceID, &device.AppVersion, &created_at, &updated_at)
		if err != nil {
			if err == sql.ErrNoRows {
				loger.Error("Failed to get device from the database", err)
			}

			return
		}

		devices.Devices = append(devices.Devices, device)
	}

	devices.Email = email

	return
}
