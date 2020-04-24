package dm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/outer-labs/forge-api-go-client/oauth"
)

const (
	DefaultBucketAPIPath = "/oss/v2/buckets"
)

// BucketAPI holds the necessary data for making Bucket related calls to Forge Data Management service
type BucketAPI struct {
	oauth.ForgeAuthenticator
	APIPath string
}

func (api BucketAPI) Path() string {
	if api.APIPath != "" {
		return api.ForgeAuthenticator.HostPath(api.APIPath)
	}
	return api.ForgeAuthenticator.HostPath(DefaultBucketAPIPath)
}

// NewBucketAPIWithCredentials returns a Bucket API client with default configurations
func NewBucketAPIWithCredentials(ClientID string, ClientSecret string) BucketAPI {
	return BucketAPI{
		ForgeAuthenticator: oauth.NewTwoLeggedClient(ClientID, ClientSecret),
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

// ErrorResult reflects the body content when a request failed (g.e. Bad request or key conflict)
type ErrorResult struct {
	Reason     string `json:"reason"`
	StatusCode int
}

func (e *ErrorResult) Error() string {
	return fmt.Sprintf("[%d]%s", e.StatusCode, e.Reason)
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
	bearer, err := api.GetTokenWithScope(oauth.ScopeBucketCreate)
	if err != nil {
		return result, err
	}
	return createBucket(api.Path(), bucketKey, policyKey, bearer.AccessToken)
}

// DeleteBucket deletes bucket given its key.
// 	WARNING: The bucket delete call is undocumented.
func (api BucketAPI) DeleteBucket(bucketKey string) error {
	bearer, err := api.GetTokenWithScope(oauth.ScopeBucketCreate)
	if err != nil {
		return err
	}
	return deleteBucket(api.Path(), bucketKey, bearer.AccessToken)
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api BucketAPI) ListBuckets(region, limit, startAt string) (result ListedBuckets, err error) {
	bearer, err := api.GetTokenWithScope(oauth.ScopeBucketRead)
	if err != nil {
		return result, err
	}
	return listBuckets(api.Path(), region, limit, startAt, bearer.AccessToken)
}

// GetBucketDetails returns information associated to a bucket. See BucketDetails struct.
func (api BucketAPI) GetBucketDetails(bucketKey string) (result BucketDetails, err error) {
	bearer, err := api.GetTokenWithScope(oauth.ScopeBucketRead)
	if err != nil {
		return result, err
	}

	return getBucketDetails(api.Path(), bucketKey, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */
func getBucketDetails(path, bucketKey, token string) (result BucketDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+bucketKey+"/details",
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

func listBuckets(path, region, limit, startAt, token string) (result ListedBuckets, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path,
		nil,
	)

	if err != nil {
		return
	}

	params := req.URL.Query()
	if len(region) != 0 {
		params.Add("region", region)
	}
	if len(limit) != 0 {
		params.Add("limit", limit)
	}
	if len(startAt) != 0 {
		params.Add("startAt", startAt)
	}

	req.URL.RawQuery = params.Encode()

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

func createBucket(path, bucketKey, policyKey, token string) (result BucketDetails, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		CreateBucketRequest{
			bucketKey,
			policyKey,
		})
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST",
		path,
		bytes.NewReader(body),
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
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

func deleteBucket(path, bucketKey, token string) (err error) {
	task := http.Client{}

	req, err := http.NewRequest("DELETE",
		path+"/"+bucketKey,
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

	return
}
