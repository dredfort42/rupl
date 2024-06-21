package db

import (
	s "auth/internal/structs"
	"database/sql"
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

// Delete access and refresh tokens from the database
func DeleteTokens(id string, accessToken string) (err error) {
	var query string
	query = `
		UPDATE ` + db.tableUsers + `
		SET access_token = NULL,
			refresh_token = NULL
		WHERE id = $1 AND access_token = $2;
	`
	_, err = db.database.Exec(query, id, accessToken)
	if err != nil {
		return
	}

	var browserIndex int
	query = `
		SELECT i
		FROM unnest((SELECT user_browsers FROM ` + db.tableUsers + ` WHERE id = $1))
		WITH ORDINALITY a(user_browser, i)
		WHERE user_browser.remembered_access_token = $2;
	`
	err = db.database.QueryRow(query, id, accessToken).Scan(&browserIndex)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return
	}

	query = `
		UPDATE ` + db.tableUsers + `
		SET user_browsers = array_remove(user_browsers, user_browsers[$1])
		WHERE id = $2;
	`
	_, err = db.database.Exec(query, browserIndex, id)

	return
}
