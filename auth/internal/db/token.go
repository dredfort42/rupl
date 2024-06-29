package db

import (
	s "auth/internal/structs"
)

// IsTokenExist checks if token exists in the database
func IsTokenExist(id string, token string, tokenType s.TokenType) (result bool) {
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
	case s.DeviceAccessToken:
		query = `
			SELECT 1
			FROM ` + db.tableDevices + `
			WHERE device_uuid = $1 AND device_access_token = $2;
		`
	case s.DeviceRefreshToken:
		query = `
			SELECT 1
			FROM ` + db.tableDevices + `
			WHERE device_uuid = $1 AND device_refresh_token = $2;
		`
	default:
		return
	}

	db.database.QueryRow(query, id, token).Scan(&result)

	return
}

// IsOneTimeRefreshToken checks if a refresh token is a one-time refresh token
func IsOneTimeRefreshToken(refreshToken string) (result bool) {
	query := `
		SELECT 1
		FROM ` + db.tableSessions + `
		WHERE refresh_token = $1 AND is_one_time = TRUE;
	`

	db.database.QueryRow(query, refreshToken).Scan(&result)

	return
}
