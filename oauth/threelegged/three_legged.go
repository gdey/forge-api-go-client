package threelegged

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

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

// Auth struct holds data necessary for making requests in 3-legged context
type Auth struct {
	oauth.AuthData
	RedirectURI string `json:"redirect_uri,omitempty"`
	Scope       scopes.Scope
}

// Authenticator interface defines the method necessary to qualify as 3-legged authenticator
type Authenticator interface {
	Authorize(state string) (string, error)
	GetToken(code string) (oauth.Bearer, error)
	RefreshToken(refreshToken string) (*oauth.Bearer, error)
}

// NewClient returns a 3-legged authenticator with default host and authPath
// if scope is 0, then ScopeDataRead is set.
func NewClient(clientID, clientSecret, redirectURI string, scope scopes.Scope) Auth {
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
		a.Path("/authorize"),
		nil,
	)

	if err != nil {
		return "", err
	}

	query := request.URL.Query()
	query.Add("client_id", a.ClientID)
	query.Add("response_type", "code")
	query.Add("redirect_uri", a.RedirectURI)
	query.Add("scope", a.Scope.String())
	query.Add("state", state)

	request.URL.RawQuery = query.Encode()

	return request.URL.String(), nil
}

// GetToken is used to exchange the authorization code for a token and an exchange token
func (a Auth) GetToken(code string) (bearer oauth.Bearer, err error) {

	task := http.Client{}

	body := url.Values{}
	body.Add("client_id", a.ClientID)
	body.Add("client_secret", a.ClientSecret)
	body.Add("grant_type", "authorization_code")
	body.Add("code", code)
	body.Add("redirect_uri", a.RedirectURI)

	req, err := http.NewRequest("POST",
		a.Path("/gettoken"),
		bytes.NewBufferString(body.Encode()),
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := task.Do(req)

	if err != nil {
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&bearer)

	return
}

// AuthToken will return an ForgeAuthenticator for the provided code
func (a Auth) AuthToken(code string) (AuthToken, error) {
	authTkn := AuthToken{Auth: a}
	bearer, err := a.GetToken(code)
	if err != nil {
		return authTkn, err
	}
	now := time.Now()
	expiryTime := now.Add(time.Second * time.Duration(bearer.ExpiresIn))

	authTkn.Token = NewRefreshableToken(&bearer, expiryTime)
	return authTkn, nil
}

// RefreshToken is used to get a new access token by using the refresh token provided by GetToken
func (a Auth) RefreshToken(refreshToken string) (bearer *oauth.Bearer, err error) {
	bearer = new(oauth.Bearer)
	task := http.Client{}

	body := url.Values{}
	body.Add("client_id", a.ClientID)
	body.Add("client_secret", a.ClientSecret)
	body.Add("grant_type", "refresh_token")
	body.Add("refresh_token", refreshToken)
	body.Add("scope", a.Scope.String())

	req, err := http.NewRequest("POST",
		a.Path("/refreshtoken"),
		bytes.NewBufferString(body.Encode()),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := task.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return nil, err
	}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(bearer)

	return nil, err
}
