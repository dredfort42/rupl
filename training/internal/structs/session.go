package structs

// LastSession is a struct for JSON
type LastSession struct {
	Email     string `json:"email"`
	StartTime int64  `json:"session_start_time"`
	EndTime   int64  `json:"session_end_time"`
}

// JSONRouteData is a struct for JSON
type JSONRouteData struct {
	Timestamp          int64   `json:"timestamp"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	HorizontalAccuracy float64 `json:"h_acc"`
	Altitude           float64 `json:"altitude"`
	VerticalAccuracy   float64 `json:"v_acc"`
	Speed              float64 `json:"speed"`
	SpeedAccuracy      float64 `json:"spd_acc"`
	Course             float64 `json:"course"`
	CourseAccuracy     float64 `json:"crs_acc"`
}

// JSONLastSessionTypeData is a struct for JSON
type JSONLastSessionTypeData struct {
	Timestamp int64  `json:"timestamp"`
	Quantity  string `json:"quantity"`
}

// JSONLastSessionData is a struct for JSON
type JSONLastSessionData struct {
	Session             LastSession               `json:"session"`
	RouteData           []JSONRouteData           `json:"route_data"`
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

// DBRouteData is a struct for the database
type DBRouteData struct {
	Timestamp          int64
	Latitude           float64
	Longitude          float64
	HorizontalAccuracy float64
	Altitude           float64
	VerticalAccuracy   float64
	Speed              float64
	SpeedAccuracy      float64
	Course             float64
	CourseAccuracy     float64
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
	SessionUUID              string
	SessionStartTime         int64
	SessionEndTime           int64
	Email                    string
	RouteData                []DBRouteData
	StepCount                []DBSessionDataInt
	RunningPower             []DBSessionDataInt
	VerticalOscillation      []DBSessionDataFloat32
	EnergyBurned             []DBSessionDataFloat32
	HeartRate                []DBSessionDataFloat32
	StrideLength             []DBSessionDataFloat32
	GroundContactTime        []DBSessionDataInt
	Speed                    []DBSessionDataFloat32
	Distance                 []DBSessionDataFloat32
	VO2MaxMLPerMinPerKg      float32
	AvrSpeedMPerS            float32
	AvrHeartRateCountPerS    float32
	AvrPowerW                float32
	AvrVerticalOscillationCm float32
	AvrGroundContactTimeMs   float32
	AvrStrideLengthM         float32
	TotalDistanceM           float32
	TotalStepsCount          int
	TotalEnergyBurnedKcal    int
}
