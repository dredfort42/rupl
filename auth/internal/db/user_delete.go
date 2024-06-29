package db

import loger "github.com/dredfort42/tools/logprinter"

// UserDelete deletes user from the database
func UserDelete(email string) (err error) {
	query := `DELETE FROM ` + db.tableUsers + ` WHERE email = $1`

	_, err = db.database.Exec(query, email)
	if err != nil {
		loger.Error("Failed to delete user from the database", err)
	}

	return
}
