package dm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	clientapi "github.com/gdey/forge-api-go-client/api"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
	"github.com/gdey/forge-api-go-client/oauth/twolegged"
)

const (
	DefaultBucketAPIPath = "/oss/v2/buckets"

	BucketFilterKeyRegion  = "region"
	BucketFilterKeyLimit   = "limit"
	BucketFilterKeyStartAt = "startAt"
)

const (
	BucketFilterRegionUS   = BucketFilterRegion(0)
	BucketFilterRegionEMEA = BucketFilterRegion(1)
)

var BucketFilterStringsRegion = [...]string{
	"US",
	"EMEA",
}

type BucketFilterRegion uint8

func (filter BucketFilterRegion) String() string {
	if filter == BucketFilterRegionEMEA {
		return BucketFilterStringsRegion[BucketFilterRegionEMEA]
	}
	return BucketFilterStringsRegion[BucketFilterRegionUS]
}

func (filter BucketFilterRegion) Add(values url.Values) (err error) {
	values.Add(BucketFilterKeyRegion, filter.String())
	return nil
}

type BucketFilterLimit uint8

func (filter BucketFilterLimit) Add(values url.Values) (err error) {
	if filter == 0 {
		filter = 10
	}
	if filter > 100 {
		filter = 100
	}
	values.Add(BucketFilterKeyLimit, strconv.FormatUint(uint64(filter), 10))
	return nil
}

type BucketFilterStartAt string

func (filter BucketFilterStartAt) Add(values url.Values) (err error) {
	if filter == "" {
		return nil
	}
	values.Add(BucketFilterKeyStartAt, string(filter))
	return nil
}

type ListBucketsFilters struct {
	Region  BucketFilterRegion
	Limit   BucketFilterLimit
	StartAt BucketFilterStartAt
}

func (filter *ListBucketsFilters) Add(values url.Values) (err error) {
	if filter == nil {
		return nil
	}
	if err = filter.Region.Add(values); err != nil {
		return err
	}
	if err = filter.Limit.Add(values); err != nil {
		return err
	}
	if err = filter.StartAt.Add(values); err != nil {
		return err
	}
	return nil
}

// BucketAPI holds the necessary data for making Bucket related calls to Forge Data Management service
type BucketAPI struct {
	Client  *clientapi.Client
	APIPath string
}

func (api BucketAPI) Path(paths ...string) []string {
	if api.APIPath == "" {
		return append([]string{DefaultBucketAPIPath}, paths...)
	}
	return append([]string{api.APIPath}, paths...)
}

// NewBucketAPIWithCredentials returns a Bucket API client with default configurations
func NewBucketAPIWithCredentials(ClientID string, ClientSecret string) BucketAPI {
	auth := twolegged.NewAuth(ClientID, ClientSecret)
	return BucketAPI{
		Client: clientapi.NewClient(auth),
	}
}

// CreateBucketRequest contains the data necessary to be passed upon bucket creation
type CreateBucketRequest struct {
	BucketKey string `json:"bucketKey"`
	PolicyKey string `json:"policyKey"`
}

// BucketDetails reflects the body content received upon creation of a bucket
type BucketDetails struct {
	BucketKey   string `json:"bucketKey"`
	BucketOwner string `json:"bucketOwner"`
	CreateDate  string `json:"createDate"`
	Permissions []struct {
		AuthID string `json:"authId"`
		Access string `json:"access"`
	} `json:"permissions"`
	PolicyKey string `json:"policyKey"`
}

// ListedBuckets reflects the response when query Data Management API for buckets associated with current Forge secrets.
type ListedBuckets struct {
	Items []struct {
		BucketKey   string `json:"bucketKey"`
		CreatedDate uint64 `json:"createdDate"`
		PolicyKey   string `json:"policyKey"`
	} `json:"items"`
	Next string `json:"next"`
}

// CreateBucket creates and returns details of created bucket, or an error on failure
func (api BucketAPI) CreateBucket(bucketKey, policyKey string) (result BucketDetails, err error) {

	body, err := json.Marshal(
		CreateBucketRequest{
			bucketKey,
			policyKey,
		})
	if err != nil {
		return result, err
	}
	err = api.Client.Post(
		context.Background(),
		scopes.BucketRead,
		api.Path(),
		&result,
		clientapi.ContentTypeJSON,
		bytes.NewReader(body),
	)
	return result, err
}

// DeleteBucket deletes bucket given its key.
// 	WARNING: The bucket delete call is undocumented.
func (api BucketAPI) DeleteBucket(bucketKey string) error {
	return api.Client.Delete(
		context.Background(),
		scopes.BucketRead,
		api.Path(bucketKey),
	)
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api BucketAPI) ListBuckets(filters *ListBucketsFilters) (result ListedBuckets, err error) {
	err = api.Client.Get(
		context.Background(),
		scopes.BucketRead,
		api.Path(),
		&result,
		filters,
	)
	return result, err
}

// GetBucketDetails returns information associated to a bucket. See BucketDetails struct.
func (api BucketAPI) GetBucketDetails(bucketKey string) (result BucketDetails, err error) {
	err = api.Client.Get(
		context.Background(),
		scopes.BucketRead,
		api.Path(bucketKey, "details"),
		&result,
	)
	return result, err
}
