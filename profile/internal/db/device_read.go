package db

// GetDevices returns a device from the database
func GetDevices(email string) ([]Device, error) {
	var devices []Device

	query := `SELECT * FROM ` + db.tableDevices + ` WHERE email = $1;`

	rows, err := db.database.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var device Device
		if err := rows.Scan(&device.DeviceModel, &device.DeviceName, &device.SystemName, &device.SystemVersion, &device.DeviceID, &device.AppVersion); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}
