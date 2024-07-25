package db

import (
	"database/sql"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// DeviceExistsCheck checks if a device exists in the database based on the email and device ID provided
func DeviceExistsCheck(email string, deviceUUID string) (result bool) {
	query := `
		SELECT 1 
		FROM ` + db.tableDevices + ` 
		WHERE email = $1
		AND device_uuid = $2;
	`

	err := db.database.QueryRow(query, email, deviceUUID).Scan(&result)
	if err != nil && err != sql.ErrNoRows {
		loger.Error("Failed to check if device exists in the database", err)
	}

	return
}

// UserDevicesExistsCheck checks if a user has any devices in the database based on the email provided
func UserDevicesExistsCheck(email string) (result bool) {
	query := `
		SELECT 1 
		FROM ` + db.tableDevices + ` 
		WHERE email = $1;
	`

	err := db.database.QueryRow(query, email).Scan(&result)
	if err != nil && err != sql.ErrNoRows {
		loger.Error("Failed to check if user has any devices in the database", err)
	}

	return
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
		var tmpEmail string
		var device s.Device
		var created_at string
		var updated_at string

		err = rows.Scan(&tmpEmail, &device.DeviceUUID, &device.DeviceModel, &device.DeviceName, &device.SystemName, &device.SystemVersion, &device.AppVersion, &created_at, &updated_at)
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
