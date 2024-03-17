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
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}
