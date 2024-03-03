package db

import (
	"github.com/dredfort42/tools/logprinter"
)

// CheckUserExists checks if a user exists in the database
func CheckUserExists(email string) bool {
	query := `SELECT email FROM ` + db.tableUsers + ` WHERE email = $1`

	if err := db.database.QueryRow(query, email).Scan(&email); err != nil {
		logprinter.PrintError("Failed to check if user exists in the database", err)
		return false
	}
	return true
}

// CheckUserPassword checks if a user's password is correct
func CheckUserPassword(email string, password string) bool {
	query := `SELECT email FROM ` + db.tableUsers + ` WHERE email = $1 AND password_hash = crypt($2, password_hash)`
	if err := db.database.QueryRow(query, email, password).Scan(&email); err != nil {
		logprinter.PrintError("Failed to check if user password is correct", err)
		return false
	}
	return true
}

// CheckUserEmailVerified checks if a user's email_verified status is true
func CheckUserEmailVerified(email string) bool {
	query := `SELECT email FROM ` + db.tableUsers + ` WHERE email = $1 AND email_verified = TRUE`
	if err := db.database.QueryRow(query, email).Scan(&email); err != nil {
		logprinter.PrintError("Failed to check if user email_verified status is true", err)
		return false
	}
	return true
}

// CheckUserAccessToken checks if a user's access token is correct
func CheckUserAccessToken(email string, accessToken string) bool {
	query := `SELECT email FROM ` + db.tableUsers + ` WHERE email = $1 AND access_token = $2`
	if err := db.database.QueryRow(query, email, accessToken).Scan(&email); err != nil {
		logprinter.PrintError("Failed to check if user access token is correct", err)
		return false
	}
	return true
}

// CheckUserRefreshToken checks if a user's refresh token is correct
func CheckUserRefreshToken(email string, refreshToken string) bool {
	query := `SELECT email FROM ` + db.tableUsers + ` WHERE email = $1 AND refresh_token = $2`
	if err := db.database.QueryRow(query, email, refreshToken).Scan(&email); err != nil {
		logprinter.PrintError("Failed to check if user refresh token is correct", err)
		return false
	}
	return true
}
