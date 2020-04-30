package dm

import (
	"context"
	"net/url"
	"strconv"

	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

const (
	ProjectFilterKeyID        = "filter[id]"
	ProjectFilterKeyType      = "filter[extension.type]"
	ProjectFilterKeyPageNum   = "page[number]"
	ProjectFilterKeyPageLimit = "page[limit]"
)

type ProjectFilterID []string

func (filter ProjectFilterID) Add(values url.Values) error {
	for _, val := range filter {
		values.Add(ProjectFilterKeyID, val)
	}
	return nil
}

type ProjectFilterType []string

func (filter ProjectFilterType) Add(values url.Values) error {
	for _, val := range filter {
		values.Add(ProjectFilterKeyType, val)
	}
	return nil
}

type ProjectFilterPageNumber int

func (filter ProjectFilterPageNumber) Add(values url.Values) error {
	if filter == 0 {
		return nil
	}
	values.Add(ProjectFilterKeyPageNum, strconv.FormatInt(int64(filter), 10))
	return nil
}

type ProjectFilterPageLimit int

func (filter ProjectFilterPageLimit) Add(values url.Values) error {
	if filter == 0 {
		return nil
	}
	values.Add(ProjectFilterKeyPageLimit, strconv.FormatInt(int64(filter), 10))
	return nil
}

type ListProjectFilters struct {
	Number ProjectFilterPageNumber
	Limit  ProjectFilterPageLimit
	ID     ProjectFilterID
	Type   ProjectFilterType
}

func (filter *ListProjectFilters) Add(values url.Values) (err error) {
	if filter == nil {
		return nil
	}
	if err = filter.ID.Add(values); err != nil {
		return err
	}
	if err = filter.Type.Add(values); err != nil {
		return err
	}
	if err = filter.Number.Add(values); err != nil {
		return err
	}
	if err = filter.Limit.Add(values); err != nil {
		return err
	}
	return nil
}

// ListProjects returns a list of all buckets created or associated with Forge secrets used for token creation
func (api HubAPI) ListProjects(hubKey string, filters *ListProjectFilters) (result ForgeResponseArray, err error) {

	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(hubKey, "projects"),
		&result,
		filters,
	)
	return result, err
}

func (api HubAPI) GetProjectDetails(hubKey, projectKey string) (result ForgeResponseObject, err error) {

	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(hubKey, "projects", projectKey),
		&result,
	)
	return result, err
}

func (api HubAPI) GetTopFolders(hubKey, projectKey string) (result ForgeResponseArray, err error) {
	err = api.Client.Get(
		context.Background(),
		scopes.DataRead,
		api.Path(hubKey, "projects", projectKey, "topFolders"),
		&result,
	)
	return result, err
}
