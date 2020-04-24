package dm

import (
	// "fmt"
	"encoding/json"
	"net/http"

	"github.com/gdey/forge-api-go-client/oauth"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
	"github.com/gdey/forge-api-go-client/oauth/twolegged"
)

const (
	DefaultHubAPIPAth = "/project/v1/hubs"
)

// HubAPI holds the necessary data for making calls to Forge Data Management service
type HubAPI struct {
	oauth.ForgeAuthenticator
	APIPath string
}

var api HubAPI

// NewHubAPIWithCredentials returns a Hub API client with default configurations
func NewHubAPIWithCredentials(ClientID string, ClientSecret string) HubAPI {
	return HubAPI{
		ForgeAuthenticator: twolegged.NewClient(ClientID, ClientSecret),
	}
}

func (api HubAPI) Path() string {
	if api.APIPath != "" {
		return api.ForgeAuthenticator.HostPath(api.APIPath)
	}
	return api.ForgeAuthenticator.HostPath(DefaultHubAPIPAth)
}

func (api HubAPI) GetHubs() (result ForgeResponseArray, err error) {
	bearer, err := api.GetTokenWithScope(scopes.DataRead)
	if err != nil {
		return
	}
	return getHubs(api.Path(), bearer.AccessToken)
}

func (api HubAPI) GetHubDetails(hubKey string) (result ForgeResponseObject, err error) {
	bearer, err := api.GetTokenWithScope(scopes.DataRead)
	if err != nil {
		return
	}
	return getHubDetails(api.Path(), hubKey, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */

func getHubs(path, token string) (result ForgeResponseArray, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path,
		nil,
	)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if response.StatusCode != http.StatusOK {
		err = &ErrorResult{StatusCode: response.StatusCode}
		decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}

func getHubDetails(path, hubKey, token string) (result ForgeResponseObject, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+hubKey,
		nil,
	)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if response.StatusCode != http.StatusOK {
		err = &ErrorResult{StatusCode: response.StatusCode}
		decoder.Decode(err)
		return
	}

	err = decoder.Decode(&result)

	return
}
