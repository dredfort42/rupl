package structs

// DeviceAuthorizationResponse is a struct for JSON
type DeviceAuthorizationResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// DeviceAccessTokenResponse is a struct for JSON
type DeviceAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// UserCredentials is a struct for JSON
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUserRequest is a struct for JSON
type LoginUserRequest struct {
	UserCredentials
	Remember bool `json:"remember"`
}

// ResponseError is a struct for JSON error
type ResponseError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
