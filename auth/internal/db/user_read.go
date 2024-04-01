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

// CheckUserRememberMe checks if a user's remember_me status is true
func CheckUserRememberMe(email string) bool {
	query := `SELECT email FROM ` + db.tableUsers + ` WHERE email = $1 AND remember_me = TRUE`
	if err := db.database.QueryRow(query, email).Scan(&email); err != nil {
		logprinter.PrintError("Failed to check if user remember_me status is true", err)
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

// CheckDeviceAccessToken checks if a device's access token is correct
func CheckDeviceAccessToken(clientID string, accessToken string) bool {
	query := `SELECT device_uuid FROM ` + db.tableUsers + ` WHERE device_uuid = $1 AND device_access_token = $2`
	if err := db.database.QueryRow(query, clientID, accessToken).Scan(&clientID); err != nil {
		logprinter.PrintError("Failed to check if device access token is correct", err)
		return false
	}
	return true
}

// GetEmailByAccessToken returns the email based on the access token provided
func GetEmailByAccessToken(accessToken string) string {
	var email string

	query := `SELECT email FROM ` + db.tableUsers + ` WHERE access_token = $1`

	if err := db.database.QueryRow(query, accessToken).Scan(&email); err != nil {
		query = `SELECT email FROM ` + db.tableUsers + ` WHERE device_access_token = $1`

		if err := db.database.QueryRow(query, accessToken).Scan(&email); err != nil {
			logprinter.PrintError("Failed to get email by access token from the database", err)
			return ""
		}
	}

	return email
}
