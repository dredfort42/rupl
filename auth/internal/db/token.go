package db

import (
	s "auth/internal/structs"
)

// DoesTokenExist checks if token exists in the database
func DoesTokenExist(id string, token string, tokenType s.TokenType) (result bool) {
	var query string

	switch tokenType {
	case s.AccessToken:
		query = `
			SELECT 1
			FROM ` + db.tableSessions + `
			WHERE email = $1 AND access_token = $2;
		`
	case s.RefreshToken:
		query = `
			SELECT 1
			FROM ` + db.tableSessions + `
			WHERE email = $1 AND refresh_token = $2;
		`
	case s.DeviceToken:
		query = `
			SELECT 1
			FROM ` + db.tableDevices + `
			WHERE email = $1 AND device_token = $2;
		`
	default:
		return
	}

	db.database.QueryRow(query, id, token).Scan(&result)

	return
}

// DoesOneTimeRefreshToken checks if a refresh token is a one-time refresh token
func DoesOneTimeRefreshToken(refreshToken string) (result bool) {
	query := `
		SELECT 1
		FROM ` + db.tableSessions + `
		WHERE refresh_token = $1 AND is_one_time = TRUE;
	`

	db.database.QueryRow(query, refreshToken).Scan(&result)

	return
}

// // IsRefreshTokenRemembered checks if a user's refresh token is remembered in the user browsers
// func IsRefreshTokenRemembered(id string, refreshToken string) (result bool) {
// 	query := `
// 		SELECT 1
// 		FROM unnest((SELECT user_browsers FROM ` + db.tableUsers + ` WHERE id = $1)) AS ub
// 		WHERE ub.remembered_refresh_token = $2;
// 	`
// 	err := db.database.QueryRow(query, id, refreshToken).Scan(&result)
// 	if err != nil {
// 		loger.Debug(err.Error())
// 	}

// 	return
// }
