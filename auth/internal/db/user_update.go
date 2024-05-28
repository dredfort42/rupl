package db

// UpdateUserEmail updates a user's email in the database
func UpdateUserEmail(oldEmail string, newEmail string) (err error) {
	query := `
		UPDATE ` + db.tableUsers + `
		SET 
			email = $2,
			updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1
	`

	_, err = db.database.Exec(query, oldEmail, newEmail)

	return
}

// UpdateUserPassword updates a user's password in the database
func UpdateUserPassword(email string, password string) (err error) {
	query := `
		UPDATE ` + db.tableUsers + ` 
		SET 
			password_hash = crypt($2, gen_salt('bf')), 
			updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1
	`

	_, err = db.database.Exec(query, email, password)

	return
}

// UpdateUserRememberMe updates a user's remember_me status in the database
func UpdateUserRememberMe(email string, rememberMe bool) (err error) {
	query := `
		UPDATE ` + db.tableUsers + `
		SET 
			remember_me = $2, 
			updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1
	`

	_, err = db.database.Exec(query, email, rememberMe)

	return
}

// UpdateUserTokens updates a user's tokens in the database
func UpdateUserTokens(email string, accessToken string, refreshToken string) (err error) {
	query := `
		UPDATE ` + db.tableUsers + ` 
		SET 
			access_token = $2, 
			refresh_token = $3, 
			updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1
	`

	_, err = db.database.Exec(query, email, accessToken, refreshToken)

	return
}

// UpdateUserEmailVerified updates a user's email_verified status in the database
func UpdateUserEmailVerified(email string) (err error) {
	query := `
		UPDATE ` + db.tableUsers + ` 
		SET 
			email_verified = TRUE, 
			updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1
	`

	_, err = db.database.Exec(query, email)

	return
}

// ================================================================================================================
// // UpdateUserDeviceUUID updates a user's device_uuid in the database
// func UpdateUserDeviceUUID(email string, deviceUUID string) {
// 	query := `UPDATE ` + db.tableUsers + ` SET
// 		device_uuid = $2,
// 		updated_at = CURRENT_TIMESTAMP
// 		WHERE email = $1`

// 	if _, err = db.database.Exec(query, email, deviceUUID); err != nil {
// 		loger.Error("Failed to update user device_uuid in the database", err)
// 	} else {
// 		loger.Debug("User device_uuid successfully updated in the database")
// 	}
// }

// // UpdateUserDeviceAccessToken updates a user's device_access_token in the database
// func UpdateUserDeviceAccessToken(email string, deviceAccessToken string) {
// 	query := `UPDATE ` + db.tableUsers + ` SET
// 		device_access_token = $2,
// 		updated_at = CURRENT_TIMESTAMP
// 		WHERE email = $1`

// 	if _, err = db.database.Exec(query, email, deviceAccessToken); err != nil {
// 		loger.Error("Failed to update user device_access_token in the database", err)
// 	} else {
// 		loger.Debug("User device_access_token successfully updated in the database")
// 	}
// }
