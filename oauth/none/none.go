package none

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gdey/forge-api-go-client/oauth"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

var None = Auth{}

type Auth struct{}

func (_ Auth) GetTokenWithScope(scope scopes.Scope) (*oauth.Bearer, error) {
	return nil, errors.New("unsupported")
}
func (_ Auth) Path(paths ...string) string                                { return strings.Join(paths, "/") }
func (_ Auth) SetAuthHeader(scope scopes.Scope, header http.Header) error { return nil }

func (_ Auth) HostPath(rest string) string { return rest }
