package db

import loger "github.com/dredfort42/tools/logprinter"

// // DeviceDelete deletes a device from the database
// func DeviceDelete(email string, deviceUUID string) (err error) {
// 	query := `
// 		DELETE FROM ` + db.tableDevices + `
// 		WHERE email = $1
// 		AND device_id = $2;
// 	`

// 	_, err = db.database.Exec(query, email, deviceUUID)
// 	if err != nil {
// 		loger.Error("Failed to delete device from the database", err)
// 	}

// 	return
// }

// DeviceDeleteAll deletes all devices from the database
func DeviceDeleteAll(email string) (err error) {
	query := `
		DELETE FROM ` + db.tableDevices + `
		WHERE email = $1;
	`

	_, err = db.database.Exec(query, email)
	if err != nil {
		loger.Error("Failed to delete all devices from the database", err)
	}

	return
}
