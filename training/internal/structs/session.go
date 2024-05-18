package structs

// LastSession is a struct for JSON
type LastSession struct {
	Email     string `json:"email"`
	StartTime int64  `json:"session_start_time"`
	EndTime   int64  `json:"session_end_time"`
}

// RouteData is a struct for JSON and database
type RouteData struct {
	Timestamp          int64   `json:"timestamp"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	HorizontalAccuracy float64 `json:"horizontal_accuracy"`
	Altitude           float64 `json:"altitude"`
	VerticalAccuracy   float64 `json:"vertical_accuracy"`
	Speed              float64 `json:"speed"`
	SpeedAccuracy      float64 `json:"speed_accuracy"`
	Course             float64 `json:"course"`
	CourseAccuracy     float64 `json:"course_accuracy"`
}

// JSONLastSessionTypeData is a struct for JSON
type JSONLastSessionTypeData struct {
	Timestamp int64  `json:"timestamp"`
	Quantity  string `json:"quantity"`
}

// JSONLastSessionData is a struct for JSON
type JSONLastSessionData struct {
	Session             LastSession               `json:"session"`
	RouteData           []RouteData               `json:"route_data"`
	StepCount           []JSONLastSessionTypeData `json:"step_count"`
	RunningPower        []JSONLastSessionTypeData `json:"running_power"`
	VerticalOscillation []JSONLastSessionTypeData `json:"vertical_oscillation"`
	EnergyBurned        []JSONLastSessionTypeData `json:"energy_burned"`
	HeartRate           []JSONLastSessionTypeData `json:"heart_rate"`
	StrideLength        []JSONLastSessionTypeData `json:"stride_length"`
	GroundContactTime   []JSONLastSessionTypeData `json:"ground_contact_time"`
	Speed               []JSONLastSessionTypeData `json:"speed"`
	Distance            []JSONLastSessionTypeData `json:"distance"`
	VO2Max              []JSONLastSessionTypeData `json:"vo2_max"`
}

// DBSessionDataInt is a struct for database
type DBSessionDataInt struct {
	Timestamp int64
	Data      int
}

// DBSessionDataFloat32 is a struct for database
type DBSessionDataFloat32 struct {
	Timestamp int64
	Data      float32
}

// DBSession represents a session entry in the database
type DBSession struct {
	SessionUUID              string                 `json:"session_uuid"`
	SessionStartTime         int64                  `json:"session_start_time"`
	SessionEndTime           int64                  `json:"session_end_time"`
	Email                    string                 `json:"email"`
	RouteData                []RouteData            `json:"route_data"`
	StepCount                []DBSessionDataInt     `json:"step_count"`
	RunningPower             []DBSessionDataInt     `json:"running_power"`
	VerticalOscillation      []DBSessionDataFloat32 `json:"vertical_oscillation"`
	EnergyBurned             []DBSessionDataFloat32 `json:"energy_burned"`
	HeartRate                []DBSessionDataFloat32 `json:"heart_rate"`
	StrideLength             []DBSessionDataFloat32 `json:"stride_length"`
	GroundContactTime        []DBSessionDataInt     `json:"ground_contact_time"`
	Speed                    []DBSessionDataFloat32 `json:"speed"`
	Distance                 []DBSessionDataFloat32 `json:"distance"`
	VO2MaxMLPerMinPerKg      float32                `json:"vo2max_mL_per_min_per_kg"`
	AvrSpeedMPerS            float32                `json:"avr_speed_m_per_s"`
	AvrHeartRateCountPerS    float32                `json:"avr_heart_rate_count_per_s"`
	AvrPowerW                float32                `json:"avr_power_W"`
	AvrVerticalOscillationCm float32                `json:"avr_vertical_oscillation_cm"`
	AvrGroundContactTimeMs   float32                `json:"avr_ground_contact_time_ms"`
	AvrStrideLengthM         float32                `json:"avr_stride_length_m"`
	TotalDistanceM           float32                `json:"total_distance_m"`
	TotalStepsCount          int                    `json:"total_steps_count"`
	TotalEnergyBurnedKcal    int                    `json:"total_energy_burned_kcal"`
}
