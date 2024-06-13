package db

import (
	s "auth/internal/structs"
)

// IsUserAccessTokenExists checks if a user's access token exists in the database
func IsTokenExists(id string, token string, tokenType s.TokenType) (result bool) {
	query := `
		SELECT 1 
		FROM ` + db.tableUsers + ` 
		WHERE email = $1 AND ( `

	if tokenType == s.AccessToken {
		query += `
			access_token = $2 
			OR EXISTS ( 
				SELECT 1 
				FROM unnest(user_browsers) AS ub 
				WHERE ub.remembered_access_token = $2 
			) )`
	} else if tokenType == s.RefreshToken {
		query += `
			refresh_token = $2 
			OR EXISTS ( 
				SELECT 1 
				FROM unnest(user_browsers) AS ub 
				WHERE ub.remembered_refresh_token = $2 
			) )`
	} else if tokenType == s.DeviceToken {
		query += `
			EXISTS ( 
				SELECT 1 
				FROM unnest(user_devices) AS ud 
				WHERE ud.device_access_token = $2 
			) )`
	} else {
		return
	}

	db.database.QueryRow(query, id, token).Scan(&result)
	return
}

// // Delete access and refresh tokens from the database
// func DeleteTokens(id string) (err error) {
// 	query := `
// 		UPDATE ` + db.tableUsers + `
// 		SET access_token = NULL, refresh_token = NULL
// 		WHERE email = $1`

// 	_, err = db.database.Exec(query, id)

// 	return
// }
