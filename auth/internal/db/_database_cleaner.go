package db

import (
	s "auth/internal/structs"
	"strconv"
	"time"

	loger "github.com/dredfort42/tools/logprinter"
)

// delete expired sessions from the Database
func deleteExpiredSessions() {
	query := `
		DELETE FROM ` + db.tableSessions + ` 
		WHERE (created_at < NOW() - INTERVAL '` + strconv.Itoa(s.JWTConfig.OneTimeRefreshTokenExpiration) + ` second' AND is_one_time = TRUE)
		OR (created_at < NOW() - INTERVAL '` + strconv.Itoa(s.JWTConfig.BrowserRefreshTokenExpiration) + ` second' AND is_one_time = FALSE)
	`

	_, err := db.database.Exec(query)
	if err != nil {
		loger.Debug(err.Error())
	}

}

// deleteExpiredDevices deletes expired devices from the Database
func deleteExpiredDevices() {
	query := `
		DELETE FROM ` + db.tableDevices + ` 
		WHERE created_at < NOW() - INTERVAL '` + strconv.Itoa(s.JWTConfig.DeviceTokenExpiration) + ` second'
	`

	_, err := db.database.Exec(query)
	if err != nil {
		loger.Debug(err.Error())
	}
}

// startCleaner starts the cleaner
func startCleaner() {
	go func() {
		for {
			deleteExpiredSessions()
			deleteExpiredDevices()

			time.Sleep(1 * time.Hour)
		}
	}()
}
