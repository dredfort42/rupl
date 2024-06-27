package db

import (
	s "auth/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// DoesTokenExist checks if token exists in the database
func DoesTokenExist(id string, token string, tokenType s.TokenType) (result bool) {
	query := `
		SELECT 1 
		FROM ` + db.tableUsers + ` 
		WHERE email = $1 AND ( `

	switch tokenType {
	case s.AccessToken:
		query += `access_token = $2 OR EXISTS ( SELECT 1 FROM unnest(user_browsers) AS ub WHERE ub.remembered_access_token = $2 ) );`
	case s.RefreshToken:
		query += `refresh_token = $2 OR EXISTS ( SELECT 1 FROM unnest(user_browsers) AS ub WHERE ub.remembered_refresh_token = $2 ) );`
	case s.DeviceToken:
		query += `EXISTS ( SELECT 1 FROM unnest(user_devices) AS ud WHERE ud.device_access_token = $2 ) );`
	default:
		return
	}

	err := db.database.QueryRow(query, id, token).Scan(&result)
	if err != nil {
		loger.Debug(err.Error())
	}

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
