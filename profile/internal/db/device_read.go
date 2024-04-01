package db

import "github.com/dredfort42/tools/logprinter"

// CheckDeviceExists checks if a device exists in the database based on the email and device ID provided
func CheckDeviceExists(deviceID string) bool {
	query := `SELECT email FROM ` + db.tableDevices + ` WHERE device_id = $1;`

	if err := db.database.QueryRow(query, deviceID).Scan(&deviceID); err != nil {
		if DEBUG {
			logprinter.PrintError("Failed to check if device exists in the database", err)
		}
		return false
	}
	return true
}

// GetDevices returns a device from the database
func GetDevices(email string) (UserDevices, error) {
	var devices UserDevices

	query := `SELECT * FROM ` + db.tableDevices + ` WHERE email = $1;`

	rows, err := db.database.Query(query, email)
	if err != nil {
		return UserDevices{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var tmpEmail string
		var device Device
		var created_at string
		var updated_at string

		if err := rows.Scan(&id, &tmpEmail, &device.DeviceModel, &device.DeviceName, &device.SystemName, &device.SystemVersion, &device.DeviceID, &device.AppVersion, &created_at, &updated_at); err != nil {
			if DEBUG {
				logprinter.PrintError("Failed to get device from the database", err)
			}
			return UserDevices{}, err
		}

		devices.Devices = append(devices.Devices, device)
	}

	devices.Email = email

	return devices, nil
}
