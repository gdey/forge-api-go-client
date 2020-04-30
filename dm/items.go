package dm

import (
	"context"
	"net/url"

	"github.com/gdey/forge-api-go-client/api/filters"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

func (api FolderAPI) GetItemDetails(projectKey, itemKey string) (result ForgeResponseObject, err error) {

	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(projectKey, "items", itemKey),
		&result,
		nil,
	)
	return result, err
}

func (api FolderAPI) GetItemTip(projectKey, itemKey string) (result ForgeResponseObject, err error) {

	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(projectKey, "items", itemKey, "tip"),
		&result,
		nil,
	)
	return result, err
}

type ItemVersionFilters struct {
	Pagination     *filters.Page
	VersionNumbers filters.VersionNumbers
	Type           filters.Type
	ID             filters.ID
	ExtensionType  filters.ExtensionType
}

func (filter *ItemVersionFilters) Add(values url.Values) error {
	if filter == nil {
		return nil
	}
	return filters.RunAll(values, filter.Pagination, filter.ID, filter.Type, filter.ExtensionType, filter.VersionNumbers)
}

// https://forge.autodesk.com/en/docs/data/v2/reference/http/projects-project_id-items-item_id-versions-GET/
func (api FolderAPI) GetItemVersions(projectKey, itemKey string, filter *ItemVersionFilters) (result ForgeResponseArray, err error) {

	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(projectKey, "items", itemKey, "versions"),
		&result,
		filter,
	)
	return result, err
}
