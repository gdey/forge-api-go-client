package dm

import (
	"context"
	"fmt"

	clientapi "github.com/gdey/forge-api-go-client/api"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
	"github.com/gdey/forge-api-go-client/oauth/twolegged"
)

const (
	DefaultFolderAPIPath = "/data/v1/projects"
)

// ErrorResult reflects the body content when a request failed (g.e. Bad request or key conflict)
type ErrorResult struct {
	Reason     string `json:"reason"`
	StatusCode int
}

func (e *ErrorResult) Error() string {
	return fmt.Sprintf("[%d]`%s`", e.StatusCode, e.Reason)
}

// FolderAPI holds the necessary data for making calls to Forge Data Management service
type FolderAPI struct {
	Client  *clientapi.Client
	APIPath string
}

// NewFolderAPIWithCredentials returns a Folder API client with default configurations
func NewFolderAPIWithCredentials(ClientID string, ClientSecret string) FolderAPI {
	auth := twolegged.NewAuth(ClientID, ClientSecret)
	return FolderAPI{
		Client: clientapi.NewClient(auth),
	}
}

func (api FolderAPI) Path(paths ...string) []string {
	if api.APIPath == "" {
		return append([]string{DefaultFolderAPIPath}, paths...)
	}
	return append([]string{api.APIPath}, paths...)
}

func (api FolderAPI) GetFolderDetails(projectKey, folderKey string) (result ForgeResponseObject, err error) {

	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(projectKey, "folders", folderKey),
		&result,
		nil,
	)
	return result, err
}

func (api FolderAPI) GetFolderContents(projectKey, folderKey string) (result ForgeResponseArray, err error) {
	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(projectKey, "folders", folderKey, "contents"),
		&result,
		nil,
	)
	return result, err
}

func (api FolderAPI) GetFolders(projectKey string) (result ForgeResponseArray, err error) {
	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(projectKey, "folders"),
		&result,
		nil,
	)
	return result, err
}
