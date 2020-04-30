package dm

import (
	// "fmt"
	"context"
	"net/url"
	"strings"

	clientapi "github.com/gdey/forge-api-go-client/api"
	"github.com/gdey/forge-api-go-client/api/filters"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
	"github.com/gdey/forge-api-go-client/oauth/twolegged"
)

const (
	DefaultHubAPIPath = "/project/v1/hubs"
)

// HubFilterType to filter by the type of the hub
type HubFilterType uint8

const (
	// HubFilterType Will filter return only hub that are part of a TEAM
	HubFilterTypeTeam = HubFilterType(1 << iota)

	// HubFilterType Will filter return only BIM360 hubs
	HubFilterTypeBIM360

	// HubFilterType Will filter return only personal hubs
	HubFilterTypePersonal

	// HubFilterTypeAll will return all hubs
	HubFilterTypeALL = HubFilterType(0)
	HubFilterKeyType = filters.KeyExtensionType
)

var hubFilterStringsType = [...]string{
	"hubs:autodesk.core:Hub",
	"hubs:autodesk.bim360:Account",
	"hubs:autodesk.a360:PersonalHub",
}

func (filter HubFilterType) Add(values url.Values) error {
	for i, val := range hubFilterStringsType {
		if (filter & (1 << i)) == (1 << i) {
			values.Add(HubFilterKeyType, val)
		}
	}
	return nil
}
func (filter HubFilterType) String() string {
	var str strings.Builder
	for i, val := range hubFilterStringsType {
		if (filter & (1 << i)) == (1 << i) {
			if i != 0 {
				str.WriteRune(' ')
			}
			str.WriteString(val)
		}
	}
	return str.String()
}

type HubsFilters struct {
	Type HubFilterType
	ID   filters.ID
	Name filters.Name
}

func (filter *HubsFilters) Add(values url.Values) (err error) {
	if filter == nil {
		return nil
	}
	return filters.RunAll(values, filter.Type, filter.ID, filter.Name)
}

// HubAPI holds the necessary data for making calls to Forge Data Management service
type HubAPI struct {
	Client  *clientapi.Client
	APIPath string
}

var api HubAPI

// NewHubAPIWithCredentials returns a Hub API client with default configurations
func NewHubAPIWithCredentials(ClientID string, ClientSecret string) HubAPI {
	return HubAPI{
		Client: clientapi.NewClient(twolegged.NewAuth(ClientID, ClientSecret)),
	}
}

func (api HubAPI) Path(paths ...string) []string {
	if api.APIPath == "" {
		return append([]string{DefaultHubAPIPath}, paths...)
	}
	return append([]string{api.APIPath}, paths...)
}

// GetHubs returns a list of know hubs
func (api HubAPI) GetHubs(hubFilters *HubsFilters) (result ForgeResponseArray, err error) {
	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(),
		&result,
		hubFilters,
	)
	return result, err
}

// GetHubDetails returns the Details for the given hub
func (api HubAPI) GetHubDetails(hubKey string) (result ForgeResponseObject, err error) {
	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(hubKey),
		&result,
	)
	return result, err
}
