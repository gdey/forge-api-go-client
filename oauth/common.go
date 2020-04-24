package oauth

import (
	"strings"
	"time"
)

const (
	// DefaultHost is the default host for autodesk api
	DefaultHost = "https://developer.api.autodesk.com"
	// DefaultAuthPath is the default AuthPath for autodesk api
	DefaultAuthPath = "/authentication/v1"
)

// Bearer reflects the response when acquiring a 2-legged token or in 3-legged context for exchanging the authorization
// code for a token + refresh token and when exchanging the refresh token for a new token
type Bearer struct {
	TokenType    string `json:"token_type"`              // Will always be Bearer
	ExpiresIn    int32  `json:"expires_in"`              // Access token expiration time (in seconds)
	AccessToken  string `json:"access_token"`            // The access token
	RefreshToken string `json:"refresh_token,omitempty"` // The refresh token used in 3-legged oauth
}

// AuthData reflects the data common to 2-legged and 3-legged api calls
type AuthData struct {
	ClientID        string    `json:"client_id,omitempty"`
	ClientSecret    string    `json:"client_secret,omitempty"`
	Host            string    `json:"host,omitempty"`
	AuthPath        string    `json:"auth_path"`
	TokenExpireTime time.Time `json:"expire_time,omitempty"` // Calculated expiration time against time.Now() for 3-legged oauth
}

// AuthDataForClient will create a new AuthData object with client info, and now for
// Expiry time
func AuthDataForClient(id, secret string) AuthData {
	return AuthData{
		ClientID:        id,
		ClientSecret:    secret,
		TokenExpireTime: time.Now(),
	}
}

// Path will return a host and auth path prepended to rest
func (a AuthData) Path(rest string) string {
	var str strings.Builder
	if a.Host != "" {
		str.WriteString(a.Host)

	} else {
		str.WriteString(DefaultHost)
	}
	if a.AuthPath != "" {
		str.WriteString(a.AuthPath)
	} else {
		str.WriteString(DefaultAuthPath)

	}
	str.WriteString(rest)
	return str.String()
}

// HostPath will return a path with the host prepended to rest
func (a AuthData) HostPath(rest string) string {
	var str strings.Builder
	if a.Host != "" {
		str.WriteString(a.Host)
	} else {
		str.WriteString(DefaultHost)
	}
	str.WriteString(rest)
	return str.String()
}

// ForgeAuthenticator defines an interface that allows abstraction from
// a 2-legged and a 3-legged context.
// 	This provides useful when an API accepts both 2-legged and 3-legged context tokens
type ForgeAuthenticator interface {
	GetTokenWithScope(scope Scopes) (*Bearer, error)
	HostPath(string) string
}
