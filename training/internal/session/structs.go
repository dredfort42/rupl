package session

// Session is a struct for JSON
type Session struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	PlanID    int    `json:"plan_id"`
	Completed bool   `json:"completed"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at"`
}
