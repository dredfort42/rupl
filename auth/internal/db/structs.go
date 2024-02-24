package db

import (
	"database/sql"
)

// Database is the database struct
type Database struct {
	database   *sql.DB
	tableUsers string
	err        error
}
