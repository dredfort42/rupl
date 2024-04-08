package db

import (
	"database/sql"
)

// Database is the database struct
type Database struct {
	database      *sql.DB
	tablePlans    string
	tableSessions string
	err           error
}
