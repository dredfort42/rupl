package db

import (
	"database/sql"
	s "profile/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
)

// UserDelete deletes a user profile and all associated devices from the database
func UserDelete(email string) (err error) {
	if !UserExistsCheck(email) {
		loger.Warning("User does not exist in the database")
		return
	}

	query := `
		DELETE FROM ` + db.tableUsers + ` 
		WHERE email = $1;
	`

	_, err = db.database.Exec(query, email)
	if err != nil {
		loger.Error("Failed to delete user from the database", err)
		return
	}

	if !UserDevicesExistsCheck(email) {
		loger.Warning("Device does not exist in the database")
		return
	}

	query = `
		DELETE FROM ` + db.tableDevices + ` 
		WHERE email = $1;
	`

	_, err = db.database.Exec(query, email)
	if err != nil {
		loger.Error("Failed to delete devices from the database", err)
	}

	return
}
