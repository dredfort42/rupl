package db

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
		var device Device
		if err := rows.Scan(&device.DeviceModel, &device.DeviceName, &device.SystemName, &device.SystemVersion, &device.DeviceID, &device.AppVersion); err != nil {
			return UserDevices{}, err
		}

		devices.Devices = append(devices.Devices, device)
	}

	devices.Email = email

	return devices, nil
}
