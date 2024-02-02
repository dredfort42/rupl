package api

// // DeviceAuthorizationRequest is a struct for JSON
// type DeviceAuthorizationRequest struct {
// 	ClientID string `json:"client_id"`
// }

// DeviceAuthorizationResponse is a struct for JSON
type DeviceAuthorizationResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// DeviceAuthorizationError is a struct for JSON
type DeviceAuthorizationError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
