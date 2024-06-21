package db

import (
	"database/sql"

	loger "github.com/dredfort42/tools/logprinter"
)

// RememberUserTokens remembers a user's tokens in the database
func RememberUserTokens(email string, accessToken string, refreshToken string) (err error) {
	query := `
		UPDATE ` + db.tableUsers + ` 
		SET 
			access_token = NULL,
			refresh_token = NULL,
			user_browsers = array_append(user_browsers, ROW($2, $3)::user_browsers),
			updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1
	`

	_, err = db.database.Exec(query, email, accessToken, refreshToken)
	return
}

// UpdateUserTokens updates a user's tokens in the database
func UpdateUserTokens(email string, newAccessToken string, newRefreshToken string) (err error) {
	query := `
		UPDATE ` + db.tableUsers + ` 
		SET 
			access_token = $2, 
			refresh_token = $3, 
			updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1
	`

	_, err = db.database.Exec(query, email, newAccessToken, newRefreshToken)
	return
}

// UpdateBrowsersTokens updates a user browser's tokens in the database
func UpdateBrowsersTokens(email string, refreshToken string, newAccessToken string, newRefreshToken string) (err error) {
	query := `
		UPDATE ` + db.tableUsers + `
		SET user_browsers = (
			SELECT array_agg(ub) FROM (
				SELECT 
					CASE 
						WHEN ub.remembered_refresh_token = $2 THEN
							ROW($3, $4)::user_browsers
						ELSE
							ub
					END
				FROM unnest(user_browsers) AS ub
			)
		)
		WHERE email = $1;
	`

	_, err = db.database.Exec(query, email, refreshToken, newAccessToken, newRefreshToken)
	return
}

// Delete access and refresh tokens from the database
func DeleteTokens(email string, accessToken string) (err error) {
	var query string
	query = `
		UPDATE ` + db.tableUsers + `
		SET access_token = NULL,
			refresh_token = NULL, 
			updated_at = CURRENT_TIMESTAMP
		WHERE email = $1 AND access_token = $2;
	`
	_, err = db.database.Exec(query, email, accessToken)
	if err != nil {
		loger.Debug(err.Error())
		return
	}

	var browserIndex int
	query = `
		SELECT i
		FROM unnest((SELECT user_browsers FROM ` + db.tableUsers + ` WHERE email = $1))
		WITH ORDINALITY a(user_browser, i)
		WHERE user_browser.remembered_access_token = $2;
	`
	err = db.database.QueryRow(query, email, accessToken).Scan(&browserIndex)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		loger.Debug(err.Error())
		return
	}

	query = `
		UPDATE ` + db.tableUsers + `
		SET user_browsers = array_remove(user_browsers, user_browsers[$1]), 
			updated_at = CURRENT_TIMESTAMP
		WHERE email = $2;
	`
	_, err = db.database.Exec(query, browserIndex, email)
	if err != nil {
		loger.Debug(err.Error())
	}

	return
}

// // UpdateUserEmail updates a user's email in the database
// func UpdateUserEmail(oldEmail string, newEmail string) (err error) {
// 	query := `
// 		UPDATE ` + db.tableUsers + `
// 		SET
// 			email = $2,
// 			updated_at = CURRENT_TIMESTAMP
// 		WHERE email = $1
// 	`

// 	_, err = db.database.Exec(query, oldEmail, newEmail)

// 	return
// }

// // UpdateUserPassword updates a user's password in the database
// func UpdateUserPassword(email string, password string) (err error) {
// 	query := `
// 		UPDATE ` + db.tableUsers + `
// 		SET
// 			password_hash = crypt($2, gen_salt('bf')),
// 			updated_at = CURRENT_TIMESTAMP
// 		WHERE email = $1
// 	`

// 	_, err = db.database.Exec(query, email, password)

// 	return
// }

// // UpdateUserEmailVerified updates a user's email_verified status in the database
// func UpdateUserEmailVerified(email string) (err error) {
// 	query := `
// 		UPDATE ` + db.tableUsers + `
// 		SET
// 			email_verified = TRUE,
// 			updated_at = CURRENT_TIMESTAMP
// 		WHERE email = $1
// 	`

// 	_, err = db.database.Exec(query, email)

// 	return
// }

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
