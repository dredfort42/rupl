package structs

import (
	"errors"
)

// Device is a struct for JSON
type Device struct {
	DeviceModel   string `json:"device_model"`
	DeviceName    string `json:"device_name"`
	SystemName    string `json:"system_name"`
	SystemVersion string `json:"system_version"`
	DeviceUUID    string `json:"device_uuid"`
	AppVersion    string `json:"app_version"`
}

// Deveces is a struct for JSON
type UserDevices struct {
	Email   string   `json:"email"`
	Devices []Device `json:"devices"`
}

// deviceStructCheck checks if the device struct is valid.
func (dvc Device) DeviceStructCheck() error {
	if dvc.DeviceModel == "" {
		return errors.New("missing device model")
	}

	if dvc.DeviceName == "" {
		return errors.New("missing device name")
	}

	if dvc.SystemName == "" {
		return errors.New("missing system name")
	}

	if dvc.SystemVersion == "" {
		return errors.New("missing system version")
	}

	if dvc.DeviceUUID == "" {
		return errors.New("missing device UUID")
	}

	if dvc.AppVersion == "" {
		return errors.New("missing app version")
	}

	return nil
}
