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

// Interval is a struct for JSON
type Interval struct {
	ID        int `json:"id"`
	Speed     int `json:"speed"`      // in m/s
	PulseZone int `json:"pulse_zone"` // 0-5
	Distance  int `json:"distance"`   // in meters
	Duration  int `json:"duration"`   // in seconds
}

// Task is a struct for JSON
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Intervals   []Interval `json:"intervals"`
}

// Plan is a struct for JSON
// type Plan struct {
// 	Email       string `json:"email"`
// 	Description string `json:"description"`
// 	Tasks       []Task `json:"tasks"`
// }

// Session is a struct for JSON
// type Session struct {
// 	ID        int    `json:"id"`
// 	Email     string `json:"email"`
// 	PlanID    int    `json:"plan_id"`
// 	Completed bool   `json:"completed"`
// 	StartedAt string `json:"started_at"`
// 	EndedAt   string `json:"ended_at"`
// }
