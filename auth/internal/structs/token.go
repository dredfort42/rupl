package structs

// TokenType is an enumeration for different types of tokens
type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
	DeviceToken
)

// String method provides the string representation of the TokenType
func (t TokenType) String() string {
	return [...]string{"AccessToken", "RefreshToken", "DeviceToken"}[t]
}

// JWTParamethers holds JWT parameters
type JWTParamethers struct {
	TokenSecret                   string
	OneTimeAccessTokenExpiration  int
	OneTimeRefreshTokenExpiration int
	BrowserAccessTokenExpiration  int
	BrowserRefreshTokenExpiration int
	DeviceTokenExpiration         int
}

// JWTConfig holds JWT configuration
var JWTConfig JWTParamethers
