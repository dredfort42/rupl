package db

// GetDevices returns a device from the database
func GetDevices(email string) (Devices, error) {
	var devices Devices

	devices.Email = email

	query := `SELECT * FROM ` + db.tableDevices + ` WHERE email = $1;`

	rows, err := db.database.Query(query, email)
	if err != nil {
		return Devices{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var device Device
		if err := rows.Scan(&device.DeviceModel, &device.DeviceName, &device.SystemName, &device.SystemVersion, &device.DeviceID, &device.AppVersion); err != nil {
			return Devices{}, err
		}

		devices.Devices = append(devices.Devices, device)
	}
	return devices, nil
}
