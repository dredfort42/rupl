package db

// DoesUserExists checks if a user exists in the database
func DoesUserExists(email string) (result bool) {
	query := `
		SELECT 1 
		FROM ` + db.tableUsers + ` 
		WHERE email = $1
	`

	db.database.QueryRow(query, email).Scan(&result)

	return
}

// DoesUserPasswordCorrect checks if a user's password is correct
func DoesUserPasswordCorrect(email string, password string) (result bool) {
	query := `
		SELECT 1
		FROM ` + db.tableUsers + `
		WHERE email = $1
		AND password_hash = crypt($2, password_hash)
	`

	db.database.QueryRow(query, email, password).Scan(&result)

	return
}

// // IsUserRememberMeSet checks if a user's remember_me status is true
// func IsUserRememberMeSet(email string) (result bool, err error) {
// 	query := `
// 		SELECT email
// 		FROM ` + db.tableUsers + `
// 		WHERE email = $1
// 		AND remember_me = TRUE
// 	`

// 	err = db.database.QueryRow(query, email).Scan(&email)
// 	if err == nil {
// 		result = true
// 	}

// 	return
// }

// // IsUserEmailVerified checks if a user's email_verified status is true
// func IsUserEmailVerified(email string) (result bool, err error) {
// 	query := `
// 		SELECT email
// 		FROM ` + db.tableUsers + `
// 		WHERE email = $1
// 		AND email_verified = TRUE
// 	`

// 	err = db.database.QueryRow(query, email).Scan(&email)
// 	if err == nil {
// 		result = true
// 	}

// 	return
// }

// // IsUserAccessTokenExists checks if a user's access token exists in the database
// func IsUserAccessTokenExists(email string, accessToken string) (result bool, err error) {
// 	query := `
// 		SELECT email
// 		FROM ` + db.tableUsers + `
// 		WHERE email = $1
// 		AND access_token = $2
// 	`

// 	err = db.database.QueryRow(query, email, accessToken).Scan(&email)
// 	if err == nil {
// 		result = true
// 	}

// 	return
// }

// // IsUserRefreshTokenExists checks if a user's refresh token exists in the database
// func IsUserRefreshTokenExists(email string, refreshToken string) (result bool, err error) {
// 	query := `
// 		SELECT email
// 		FROM ` + db.tableUsers + `
// 		WHERE email = $1
// 		AND refresh_token = $2
// 	`

// 	err = db.database.QueryRow(query, email, refreshToken).Scan(&email)
// 	if err == nil {
// 		result = true
// 	}

// 	return
// }

// ====================================================================================================
// // CheckDeviceAccessToken checks if a device's access token is correct
// func CheckDeviceAccessToken(clientID string, accessToken string) bool {
// 	query := `SELECT device_uuid FROM ` + db.tableUsers + ` WHERE device_uuid = $1 AND device_access_token = $2`

// 	if err := db.database.QueryRow(query, clientID, accessToken).Scan(&clientID); err != nil {
// 		loger.Error("Failed to check if device access token is correct", err)
// 		return false
// 	}
// 	return true
// }

// // GetEmailByAccessToken returns the email based on the access token provided
// func GetEmailByAccessToken(accessToken string) string {
// 	var email string
// 	query := `SELECT email FROM ` + db.tableUsers + ` WHERE access_token = $1`

// 	if err := db.database.QueryRow(query, accessToken).Scan(&email); err != nil {
// 		query = `SELECT email FROM ` + db.tableUsers + ` WHERE device_access_token = $1`

// 		if err := db.database.QueryRow(query, accessToken).Scan(&email); err != nil {
// 			loger.Error("Failed to get email by access token from the database", err)
// 			return ""
// 		}
// 	}

// 	return email
// }
