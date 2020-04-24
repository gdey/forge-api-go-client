package twolegged

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gdey/forge-api-go-client/oauth"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

// Auth struct holds data necessary for making requests in 2-legged context
type Auth struct {
	oauth.AuthData
}

// Authenticator interface defines the method necessary to qualify as 2-legged authenticator
type Authenticator interface {
	Authenticate(scope scopes.Scope) (*oauth.Bearer, error)
}

// NewClient returns a 2-legged authenticator with default host and authPath
func NewClient(clientID, clientSecret string) Auth {
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
	task := http.Client{}

	body := url.Values{}
	body.Add("client_id", a.ClientID)
	body.Add("client_secret", a.ClientSecret)
	body.Add("grant_type", "client_credentials")
	body.Add("scope", scope.String())

	req, err := http.NewRequest("POST",
		a.Path("/authenticate"),
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

	return bearer, nil
}