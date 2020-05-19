package oauth

import (
	"net/http"
	"strings"

	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

const (
	// DefaultHost is the default host for autodesk api
	DefaultHost = "https://developer.api.autodesk.com"
	// DefaultAuthenticationPath is the default AuthPath for autodesk api
	DefaultAuthenticationPath = "authentication/v1"

	HeaderAuthorization = "Authorization"
	HeaderXUserID       = "x-user-id"
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
	ClientID           string `json:"client_id,omitempty"`
	ClientSecret       string `json:"client_secret,omitempty"`
	Host               string `json:"host,omitempty"`
	AuthenticationPath string `json:"auth_path"`
}

// AuthDataForClient will create a new AuthData object with client info
func AuthDataForClient(id, secret string) AuthData {
	return AuthData{
		ClientID:     id,
		ClientSecret: secret,
	}
}

// Path will return a host and auth path prepended to rest
func (a AuthData) AuthPath(rest ...string) []string {
	paths := make([]string, 2, len(rest)+2)
	if a.Host != "" {
		paths[0] = a.Host
	} else {
		paths[0] = DefaultHost
	}
	if a.AuthenticationPath != "" {
		paths[1] = a.AuthenticationPath
	} else {
		paths[1] = DefaultAuthenticationPath
	}
	return append(paths, rest...)
}

func (a AuthData) Path(rest ...string) string {
	var str strings.Builder
	if a.Host != "" {
		str.WriteString(a.Host)
	} else {
		str.WriteString(DefaultHost)
	}
	for _, astr := range rest {
		str.WriteRune('/')
		str.WriteString(astr)
	}
	return str.String()
}

// ForgeAuthenticator defines an interface that allows abstraction from
// a 2-legged and a 3-legged context.
// 	This provides useful when an API accepts both 2-legged and 3-legged context tokens
type ForgeAuthenticator interface {
	//GetTokenWithScope(scope scopes.Scope) (*Bearer, error)

	// Path returns the full url path with the given compentents
	Path(...string) string

	// SetAuthHeader should set the appropriate http headers for auth
	// It should refresh any tokens it may have
	SetAuthHeader(scope scopes.Scope, header http.Header) error
}
