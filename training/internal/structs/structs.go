package structs

// ResponseError is a struct for JSON error
type ResponseError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// Session is a struct for JSON
type Session struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	PlanID    int    `json:"plan_id"`
	Completed bool   `json:"completed"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at"`
}

// Interval is a struct for JSON
type Interval struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Distance    int    `json:"distance"` // in meters
	Duration    int    `json:"duration"` // in seconds
	Speed       Speed  `json:"speed"`
	Pulse       Pulse  `json:"pulse"`
}

type Speed struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

type Pulse struct {
	Min int `json:"min"`
	Max int `json:"max"`
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
