package dm

import (
	"context"
	"io"
	"net/http"
	"net/url"

	clientapi "github.com/gdey/forge-api-go-client/api"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

// ObjectDetails reflects the data presented when uploading an object to a bucket or requesting details on object.
type ObjectDetails struct {
	BucketKey   string            `json:"bucketKey"`
	ObjectID    string            `json:"objectID"`
	ObjectKey   string            `json:"objectKey"`
	SHA1        string            `json:"sha1"`
	Size        uint64            `json:"size"`
	ContentType string            `json:"contentType, omitempty"`
	Location    string            `json:"location"`
	BlockSizes  []int64           `json:"blockSizes, omitempty"`
	Deltas      map[string]string `json:"deltas, omitempty"`
}

const (
	ObjectFilterKeyBeginsWith = "beginsWith"
)

type ObjectFilterBeginsWith string

func (filter ObjectFilterBeginsWith) Add(values url.Values) error {
	if filter == "" {
		return nil
	}
	values.Add(ObjectFilterKeyBeginsWith, string(filter))
	return nil
}

type ListObjectsFilters struct {
	Limit      BucketFilterLimit
	StartAt    BucketFilterStartAt
	BeginsWith ObjectFilterBeginsWith
}

func (filter *ListObjectsFilters) Add(values url.Values) (err error) {
	if filter == nil {
		return nil
	}
	if err = filter.BeginsWith.Add(values); err != nil {
		return err
	}
	if err = filter.StartAt.Add(values); err != nil {
		return err
	}
	if err = filter.Limit.Add(values); err != nil {
		return err
	}
	return nil
}

// BucketContent reflects the response when query Data Management API for bucket content.
type BucketContent struct {
	Items []ObjectDetails `json:"items"`
	Next  string          `json:"next"`
}

// UploadObject adds to specified bucket the given data (can originate from a multipart-form or direct file read).
// Return details on uploaded object, including the object URN. Check ObjectDetails struct.
func (api BucketAPI) UploadObject(bucketKey string, objectName string, reader io.Reader) (result ObjectDetails, err error) {

	err = api.Client.Put(
		context.Background(),
		scopes.DataWrite|scopes.DataCreate,
		api.Path(bucketKey, "objects", objectName),
		&result,
		"",
		reader,
	)
	return result, err

}

// DownloadObject returns the reader stream of the response body
// Don't forget to close it!
// https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectName-GET/
// TODO(gdey): Create DownloadObjectOptions Struct to set various Headers
func (api BucketAPI) DownloadObject(bucketKey string, objectName string) (reader io.ReadCloser, err error) {
	res, err := api.Client.DoRawRequest(
		context.Background(), "GET",
		scopes.DataRead,
		api.Path(bucketKey, "objects", objectName),
		nil, nil, "", nil,
	)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, clientapi.ErrResult{StatusCode: res.StatusCode}
	}
	return res.Body, nil
}

// ListObjects returns the bucket contains along with details on each item.
func (api BucketAPI) ListObjects(bucketKey string, filters *ListObjectsFilters) (result BucketContent, err error) {
	err = api.Client.Get(
		context.Background(),
		scopes.BucketRead,
		api.Path(bucketKey, "objects"),
		&result,
		filters,
	)
	return result, err
}
