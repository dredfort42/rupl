package structs

// Profile is a struct for JSON
type Profile struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}

// Device is a struct for JSON
type Device struct {
	DeviceModel   string `json:"device_model"`
	DeviceName    string `json:"device_name"`
	SystemName    string `json:"system_name"`
	SystemVersion string `json:"system_version"`
	DeviceID      string `json:"device_id"`
	AppVersion    string `json:"app_version"`
}

// Deveces is a struct for JSON
type UserDevices struct {
	Email   string   `json:"email"`
	Devices []Device `json:"devices"`
}
