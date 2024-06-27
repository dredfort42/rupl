package db

// // SessionIsOneTimeExists checks if a one-time session exists in the database
// func SessionIsOneTimeExists(email string) (result bool) {
// 	query := `
// 		SELECT 1
// 		FROM ` + db.tableSessions + `
// 		WHERE email = $1
// 		AND is_one_time = true
// 	`

// 	db.database.QueryRow(query, email).Scan(&result)
// 	return
// }
