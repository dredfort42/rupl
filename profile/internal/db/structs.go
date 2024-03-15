package db

import (
	"database/sql"
)

// Database is the database struct
type Database struct {
	database      *sql.DB
	tableProfiles string
	err           error
}

// Profile is a struct for JSON
type Profile struct {
	Email       string
	FirstName   string
	LastName    string
	DateOfBirth string
	Gender      string
}
