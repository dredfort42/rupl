package structs

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

// ----------------------------

// LastSession is a struct for JSON
type LastSession struct {
	Email     string `json:"email"`
	StartTime int64  `json:"session_start_time"`
	EndTime   int64  `json:"session_end_time"`
}

// RouteData is a struct for JSON
type LastSessionRouteData struct {
	Timestamp          int64   `json:"timestamp"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	HorizontalAccuracy float64 `json:"h_acc"`
	Altitude           float64 `json:"altitude"`
	VerticalAccuracy   float64 `json:"v_acc"`
	Course             float64 `json:"course"`
	CourseAccuracy     float64 `json:"crs_acc"`
	Speed              float64 `json:"speed"`
	SpeedAccuracy      float64 `json:"spd_acc"`
}

// StepCount is a struct for JSON
type LastSessionStepCount struct {
	Timestamp  int64 `json:"timestamp"`
	StepsCount int   `json:"quantity"`
}

// RunningPower is a struct for JSON
type LastSessionRunningPower struct {
	Timestamp int64 `json:"timestamp"`
	PowerW    int   `json:"quantity"`
}

// VerticalOscillation is a struct for JSON
type LastSessionVerticalOscillation struct {
	Timestamp             int64   `json:"timestamp"`
	VerticalOscillationCm float32 `json:"quantity"`
}

// EnergyBurned is a struct for JSON
type LastSessionEnergyBurned struct {
	Timestamp        int64   `json:"timestamp"`
	EnergyBurnedKcal float32 `json:"quantity"`
}

// HeartRate is a struct for JSON
type LastSessionHeartRate struct {
	Timestamp int64   `json:"timestamp"`
	HeartRate float32 `json:"quantity"`
}

// StrideLength is a struct for JSON
type LastSessionStrideLength struct {
	Timestamp    int64   `json:"timestamp"`
	StrideLenght float32 `json:"quantity"`
}

// GroundContactTime is a struct for JSON
type LastSessionGroundContactTime struct {
	Timestamp          int64 `json:"timestamp"`
	GroundContactTimeS int   `json:"quantity"`
}

// Speed is a struct for JSON
type LastSessionSpeed struct {
	Timestamp int64   `json:"timestamp"`
	Speed     float32 `json:"quantity"`
}

// Distance is a struct for JSON
type LastSessionDistance struct {
	Timestamp int64   `json:"timestamp"`
	Distance  float32 `json:"quantity"`
}

// LastSessionData is a struct for JSON
type LastSessionData struct {
	Session             LastSession                      `json:"session"`
	RouteData           []LastSessionRouteData           `json:"route_data"`
	StepCount           []LastSessionStepCount           `json:"step_count"`
	RunningPower        []LastSessionRunningPower        `json:"running_power"`
	VerticalOscillation []LastSessionVerticalOscillation `json:"vertical_oscillation"`
	EnergyBurned        []LastSessionEnergyBurned        `json:"energy_burned"`
	HeartRate           []LastSessionHeartRate           `json:"heart_rate"`
	StrideLength        []LastSessionStrideLength        `json:"stride_length"`
	GroundContactTime   []LastSessionGroundContactTime   `json:"ground_contact_time"`
	Speed               []LastSessionSpeed               `json:"speed"`
	Distance            []LastSessionDistance            `json:"distance"`
	VO2Max              float32                          `json:"vo2_max"`
}
