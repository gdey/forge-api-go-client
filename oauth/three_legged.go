package oauth

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
)

// ThreeLeggedAuth struct holds data necessary for making requests in 3-legged context
type ThreeLeggedAuth struct {
	AuthData
	RedirectURI string `json:"redirect_uri,omitempty"`
	Scope       Scopes
}

type ThreeLeggedAuthToken struct {
	ThreeLeggedAuth
	Token *RefreshableToken
}

func (a ThreeLeggedAuthToken) GetTokenWithScope(scope Scopes) (*Bearer, error) {
	if !a.ThreeLeggedAuth.Scope.Allows(scope) {
		return nil, fmt.Errorf("scopes require: '%v' have '%v'", a.ThreeLeggedAuth.Scope, scope)
	}

	if err := a.Token.RefreshTokenIfRequired(a.ThreeLeggedAuth); err != nil {
		return nil, err
	}
	return a.Token.Bearer(), nil
}

// ThreeLeggedAuthenticator interface defines the method necessary to qualify as 3-legged authenticator
type ThreeLeggedAuthenticator interface {
	Authorize(state string) (string, error)
	GetToken(code string) (Bearer, error)
	RefreshToken(refreshToken string) (*Bearer, error)
}

// NewThreeLeggedClient returns a 3-legged authenticator with default host and authPath
// if scope is 0, then ScopeDataRead is set.
func NewThreeLeggedClient(clientID, clientSecret, redirectURI string, scope Scopes) ThreeLeggedAuth {
	if scope == 0 {
		// TOOD(gdey): would ScopeViewableRead be a better things to ask for?
		scope = ScopeDataRead
	}
	return ThreeLeggedAuth{
		AuthData:    AuthDataForClient(clientID, clientSecret),
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
func (a ThreeLeggedAuth) Authorize(state string) (string, error) {

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

//GetToken is used to exchange the authorization code for a token and an exchange token
func (a ThreeLeggedAuth) GetToken(code string) (bearer Bearer, err error) {

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

// ThreeLeggedAuthToken will return an ForgeAuthenticator for the provided code
func (a ThreeLeggedAuth) ThreeLeggedAuthToken(code string) (ThreeLeggedAuthToken, error) {
	authTkn := ThreeLeggedAuthToken{
		ThreeLeggedAuth: a,
	}
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
func (a ThreeLeggedAuth) RefreshToken(refreshToken string) (bearer *Bearer, err error) {
	bearer = new(Bearer)
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
