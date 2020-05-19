package twolegged

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gdey/forge-api-go-client/api"
	"github.com/gdey/forge-api-go-client/oauth"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

// Auth struct holds data necessary for making requests in 2-legged context
type Auth struct {
	oauth.AuthData
	// UserID is the user to act on the behalf of.
	UserID string

	client api.Client
}

// Authenticator interface defines the method necessary to qualify as 2-legged authenticator
type Authenticator interface {
	Authenticate(scope scopes.Scope) (*oauth.Bearer, error)
}

// NewClient returns a 2-legged authenticator with default host and authPath
func NewAuth(clientID, clientSecret string) Auth {
	return Auth{
		AuthData: oauth.AuthDataForClient(clientID, clientSecret),
	}
}

// GetTokenWithScope will get the a token for the given scope
func (a Auth) GetTokenWithScope(scope scopes.Scope) (*oauth.Bearer, error) {
	return a.Authenticate(scope)
}

// Authenticate allows getting a token with a given scope
func (a Auth) Authenticate(scope scopes.Scope) (bearer *oauth.Bearer, err error) {

	if !scope.IsValid() {
		return nil, errors.New("Invalid scope")
	}

	bearer = new(oauth.Bearer)

	body := url.Values{
		"client_id":     []string{a.ClientID},
		"client_secret": []string{a.ClientSecret},
		"grant_type":    []string{"client_credentials"},
		"scope":         []string{scope.String()},
	}

	res, err := a.client.DoRawRequest(context.Background(), "POST", 0,
		a.AuthPath("authenticate"),
		nil, nil,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(body.Encode()),
	)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(res.Body)
		return nil, api.ErrResult{
			StatusCode: res.StatusCode,
			Reason:     string(content),
		}
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(bearer)

	return bearer, nil
}

func (a Auth) SetAuthHeader(scope scopes.Scope, header http.Header) error {

	bearer, err := a.GetTokenWithScope(scope)
	if err != nil {
		return err
	}
	header.Set(oauth.HeaderAuthorization, "Bearer "+bearer.AccessToken)
	if a.UserID != "" {
		header.Set(oauth.HeaderXUserID, a.UserID)
	}
	return nil
}
