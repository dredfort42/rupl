package training

// Interval is a struct for JSON
type Interval struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Speed       int    `json:"speed"`      // in m/s
	PulseZone   int    `json:"pulse_zone"` // 0-5
	Distance    int    `json:"distance"`   // in meters
	Duration    int    `json:"duration"`   // in seconds
}

// Task is a struct for JSON
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Compleated  bool       `json:"compleated"`
	Intervals   []Interval `json:"intervals"`
}

// Plan is a struct for JSON
type Plan struct {
	Email       string `json:"email"`
	Description string `json:"description"`
	Tasks       []Task `json:"tasks"`
}
