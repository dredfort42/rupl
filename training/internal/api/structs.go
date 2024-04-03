package api

// ResponseError is a struct for JSON error
type ResponseError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
