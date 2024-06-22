package db

// DeleteUser deletes user from the database
func DeleteUser(email string) (err error) {
	query := `DELETE FROM ` + db.tableUsers + ` WHERE email = $1`

	_, err = db.database.Exec(query, email)

	return
}

// // Delete device access token from the database
// func DeleteDeviceAccessToken(email string) (err error) {
// 	query := `UPDATE ` + db.tableUsers + ` SET device_uuid = NULL, device_access_token = NULL WHERE email = $1`

// 	_, err = db.database.Exec(query, email)
// 	if err != nil {
// 		return
// 	}

// 	loger.Debug(email+" user's device access token successfully deleted from the database", "")
// 	return true
// }
