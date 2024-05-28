package db

// Delete user's access and refresh tokens from the database
func DeleteUserTokens(email string) (err error) {
	query := `
		UPDATE ` + db.tableUsers + ` 
		SET access_token = NULL, refresh_token = NULL 
		WHERE email = $1`

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
