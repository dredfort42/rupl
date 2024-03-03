package db

import (
	"github.com/dredfort42/tools/logprinter"
)

// UpdateUserEmail updates a user's email in the database
func UpdateUserEmail(oldEmail string, newEmail string) {
	query := `UPDATE ` + db.tableUsers + ` SET 
		email = $2, 
		updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1`

	if _, db.err = db.database.Exec(query, oldEmail, newEmail); db.err != nil {
		logprinter.PrintError("Failed to update user email in the database", db.err)
	} else {
		logprinter.PrintSuccess("User email successfully updated in the database", "")
	}
}

// UpdateUserPassword updates a user's password in the database
func UpdateUserPassword(email string, password string) {
	query := `UPDATE ` + db.tableUsers + ` SET 
		password_hash = crypt($2, gen_salt('bf')), 
		updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1`

	if _, db.err = db.database.Exec(query, email, password); db.err != nil {
		logprinter.PrintError("Failed to update user password in the database", db.err)
	} else {
		logprinter.PrintSuccess("User password successfully updated in the database", "")
	}
}

// UpdateUserTokens updates a user's tokens in the database
func UpdateUserTokens(email string, accessToken string, refreshToken string) {
	query := `UPDATE ` + db.tableUsers + ` SET 
		access_token = $2, 
		refresh_token = $3, 
		updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1`

	if _, db.err = db.database.Exec(query, email, accessToken, refreshToken); db.err != nil {
		logprinter.PrintError("Failed to update user tokens in the database", db.err)
	} else {
		logprinter.PrintSuccess("User tokens successfully updated in the database", "")
	}
}

// UpdateUserEmailVerified updates a user's email_verified status in the database
func UpdateUserEmailVerified(email string) {
	query := `UPDATE ` + db.tableUsers + ` SET 
		email_verified = TRUE, 
		updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1`

	if _, db.err = db.database.Exec(query, email); db.err != nil {
		logprinter.PrintError("Failed to update user email_verified status in the database", db.err)
	} else {
		logprinter.PrintSuccess("User email_verified status successfully updated in the database", "")
	}
}

// UpdateUserDeviceUUID updates a user's device_uuid in the database
func UpdateUserDeviceUUID(email string, deviceUUID string) {
	query := `UPDATE ` + db.tableUsers + ` SET 
		device_uuid = $2, 
		updated_at = CURRENT_TIMESTAMP 
		WHERE email = $1`

	if _, db.err = db.database.Exec(query, email, deviceUUID); db.err != nil {
		logprinter.PrintError("Failed to update user device_uuid in the database", db.err)
	} else {
		logprinter.PrintSuccess("User device_uuid successfully updated in the database", "")
	}
}
