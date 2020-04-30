package threelegged

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gdey/forge-api-go-client/api"
	"github.com/gdey/forge-api-go-client/oauth"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

type AuthToken struct {
	Auth
	Token *RefreshableToken
}

func (a AuthToken) GetTokenWithScope(scope scopes.Scope) (*oauth.Bearer, error) {
	if !a.Auth.Scope.Allows(scope) {
		return nil, fmt.Errorf("scopes require: '%v' have '%v'", scope, a.Auth.Scope)
	}

	if err := a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return nil, err
	}
	return a.Token.Bearer(), nil
}

func (a AuthToken) SetAuthHeader(scope scopes.Scope, header http.Header) error {

	bearer, err := a.GetTokenWithScope(scope)
	if err != nil {
		return err
	}
	header.Set(oauth.HeaderAuthorization, "Bearer "+bearer.AccessToken)
	return nil
}

// Auth struct holds data necessary for making requests in 3-legged context
type Auth struct {
	oauth.AuthData
	RedirectURI string `json:"redirect_uri,omitempty"`
	Scope       scopes.Scope

	// Implicate will do an implicate token retrieval.
	// Do not use this unless you know what you are doing.
	Implicate bool

	client api.Client
}

// Authenticator interface defines the method necessary to qualify as 3-legged authenticator
type Authenticator interface {
	Authorize(state string) (string, error)
	GetToken(code string) (oauth.Bearer, error)
	RefreshToken(refreshToken string) (*oauth.Bearer, error)
}

// NewAuth returns a 3-legged authenticator with default host and authPath
// if scope is 0, then ScopeDataRead is set.
func NewAuth(clientID, clientSecret, redirectURI string, scope scopes.Scope) Auth {
	if scope == 0 {
		scope = scopes.DataRead
	}
	return Auth{
		AuthData:    oauth.AuthDataForClient(clientID, clientSecret),
		RedirectURI: redirectURI,
		Scope:       scope,
	}
}

// Authorize method returns an URL to redirect an end user, where it will be asked to give his consent for app to
//access the specified resources.
//
// State can be used to specify, as URL-encoded payload, some arbitrary data that the authentication flow will pass back
// verbatim in a state query parameter to the callback URL.
//	Note: You do not call this URL directly in your server code.
//	See the Get a 3-Legged Token tutorial for more information on how to use this endpoint.
func (a Auth) Authorize(state string) (string, error) {

	request, err := http.NewRequest("GET",
		strings.Join(a.AuthData.AuthPath("authorize"), "/"),
		nil,
	)

	if err != nil {
		return "", err
	}

	query := request.URL.Query()
	query.Add("client_id", a.ClientID)
	if a.Implicate {
		query.Add("response_type", "token")

	} else {
		query.Add("response_type", "code")
	}
	query.Add("redirect_uri", a.RedirectURI)
	query.Add("scope", a.Scope.String())
	query.Add("state", state)

	request.URL.RawQuery = query.Encode()

	return request.URL.String(), nil
}

// GetToken is used to exchange the authorization code for a token and an exchange token
func (a Auth) GetToken(code string) (bearer oauth.Bearer, err error) {

	body := url.Values{
		"client_id":     []string{a.ClientID},
		"client_secret": []string{a.ClientSecret},
		"grant_type":    []string{"authorization_code"},
		"code":          []string{code},
		"redirect_uri":  []string{a.RedirectURI},
	}
	res, err := a.client.DoRawRequest(context.Background(), http.MethodPost, 0,
		a.AuthPath("gettoken"),
		nil, nil,
		api.ContentTypeFormEncoded,
		bytes.NewBufferString(body.Encode()),
	)
	if err != nil {
		return bearer, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(res.Body)
		return bearer, api.ErrorResult{
			StatusCode: res.StatusCode,
			Reason:     string(content),
		}
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&bearer)

	return bearer, nil
}

// AuthToken will return an ForgeAuthenticator for the provided code
func (a Auth) AuthToken(code string) (AuthToken, error) {
	authTkn := AuthToken{Auth: a}
	bearer, err := a.GetToken(code)
	if err != nil {
		return authTkn, err
	}

	authTkn.Token = NewRefreshableToken(&bearer)
	return authTkn, nil
}

// RefreshToken is used to get a new access token by using the refresh token provided by GetToken
func (a Auth) RefreshToken(refreshToken string) (bearer *oauth.Bearer, err error) {
	bearer = new(oauth.Bearer)

	body := url.Values{
		"client_id":     []string{a.ClientID},
		"client_secret": []string{a.ClientSecret},
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refreshToken},
		"scope":         []string{a.Scope.String()},
	}

	res, err := a.client.DoRawRequest(context.Background(), http.MethodPost, 0,
		a.AuthPath("refreshtoken"),
		nil, nil,
		api.ContentTypeFormEncoded,
		bytes.NewBufferString(body.Encode()),
	)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(res.Body)
		return nil, api.ErrorResult{
			StatusCode: res.StatusCode,
			Reason:     string(content),
		}
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(bearer)

	return bearer, nil
}
