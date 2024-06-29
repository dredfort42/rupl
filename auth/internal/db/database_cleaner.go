package db

import (
	s "auth/internal/structs"
	"strconv"
	"time"

	cfg "github.com/dredfort42/tools/configreader"
	loger "github.com/dredfort42/tools/logprinter"
)

// databaseCleanerStart starts the cleaners
func databaseCleanerStart() {
	cleanUpInrerval, err := strconv.Atoi(cfg.Config["db.cleanup.interval"])
	if err != nil {
		panic("Database cleanup interval is not set")
	}

	go func() {
		for {
			query := `
				DELETE FROM ` + db.tableSessions + `
				WHERE is_one_time = TRUE
				AND created_at < CURRENT_TIMESTAMP - INTERVAL '` + strconv.Itoa(s.JWTConfig.OneTimeAccessTokenExpiration) + ` seconds';
			`

			_, err := db.database.Exec(query)
			if err != nil {
				loger.Error("Failed to delete sessions with expired one-time tokens", err)
			}

			query = `
				DELETE FROM ` + db.tableSessions + `
				WHERE is_one_time = FALSE
				AND created_at < CURRENT_TIMESTAMP - INTERVAL '` + strconv.Itoa(s.JWTConfig.BrowserRefreshTokenExpiration) + ` seconds';
			`
			_, err = db.database.Exec(query)
			if err != nil {
				loger.Error("Failed to delete sessions with expired browser tokens", err)
			}

			time.Sleep(time.Duration(cleanUpInrerval) * time.Second)
		}
	}()
}
