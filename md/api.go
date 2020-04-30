package md

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	clientapi "github.com/gdey/forge-api-go-client/api"
	"github.com/gdey/forge-api-go-client/api/filters"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
	"github.com/gdey/forge-api-go-client/oauth/twolegged"
)

var (
	// TranslationSVFPreset specifies the minimum necessary for translating a generic (single file, uncompressed)
	// model into svf.
	TranslationSVFPreset = TranslationParams{
		Output: OutputSpec{
			Destination: DestSpec{"us"},
			Formats: []FormatSpec{
				FormatSpec{
					"svf",
					[]string{"2d", "3d"},
				},
			},
		},
	}
)

const (
	DefaultModelDerivativePath = "/modelderivative/v2/designdata"
)

//TranslationParams is used when specifying the translation jobs
type TranslationParams struct {
	Input  TranslationInput `json:"input"`
	Output OutputSpec       `json:"output"`
}
type TranslationInput struct {
	URN           string  `json:"urn"`
	CompressedURN *bool   `json:"compressedUrn,omitempty"`
	RootFileName  *string `json:"rootFileName,omitempty"`
}

// TranslationResult reflects data received upon successful creation of translation job
type TranslationResult struct {
	Result       string `json:"result"`
	URN          string `json:"urn"`
	AcceptedJobs struct {
		Output OutputSpec `json:"output"`
	}
}

type ManifestResult struct {
	Type         string           `json:"type,omitempty"`
	HasThumbnail bool             `json:"hasThumbnail,string,omitempty"`
	Status       string           `json:"status,omitempty"`
	Progress     string           `json:"progress,omitempty"`
	Region       string           `json:"region,omitempty"`
	URN          string           `json:"urn,omitempty"`
	Derivatives  []DerivativeSpec `json:"derivatives,omitempty"`
}

type DerivativeSpec struct {
	Name         string         `json:"name,omitempty"`
	HasThumbnail bool           `json:"hasThumbnail,string,omitempty"`
	Role         string         `json:"role,omitempty"`
	Status       string         `json:"status,omitempty"`
	Progress     string         `json:"progress,omitempty"`
	Children     []ChildrenSpec `json:"children,omitempty"`
}

type ChildrenSpec struct {
	GUID     string `json:"guid,omitempty"`
	Role     string `json:"role,omitempty"`
	MIME     string `json:"mime,omitempty"`
	URN      string `json:"urn,omitempty"`
	Progress string `json:"progress,omitempty"`
	Status   string `json:"status,omitempty"`
}

// OutputSpec reflects data found upon creation translation job and receiving translation job status
type OutputSpec struct {
	Destination DestSpec     `json:"destination,omitempty"`
	Formats     []FormatSpec `json:"formats"`
}

// DestSpec is used within OutputSpecs and is useful when specifying the region for translation results
type DestSpec struct {
	Region string `json:"region,omitempty"`
}

// FormatSpec is used within OutputSpecs and should be used when specifying the expected format and views (2d or/and 3d)
type FormatSpec struct {
	Type  string   `json:"type"`
	Views []string `json:"views,omitempty"`
}

type MetadataResult struct {
	Data MetadataSpec `json:"data",omitempty`
}

type MetadataSpec struct {
	Type     string     `json:"type",omitempty`
	Metadata []ViewSpec `json:"metadata",omitempty`
}

type ViewSpec struct {
	Name string `json:"name",omitempty`
	Role string `json:"role",omitempty`
	Guid string `json:"guid",omitempty`
}

type PropertiesResult struct {
	Data   PropertiesSpec `json:"data",omitempty`
	Result string         `json:"result",omitempty`
}

type PropertiesSpec struct {
	Type       string       `json:"type"`
	Collection []ObjectSpec `json:"collection"`
}

type ObjectSpec struct {
	ObjectID   int64  `json:"objectid"`
	Name       string `json:"name"`
	ExternalID string `json:"externalId"`
	Properties json.RawMessage
}

type TreeResult struct {
	Data TreeSpec `json:"data",omitempty`
}

type TreeSpec struct {
	Type    string         `json:"type",omitempty`
	Objects []TreeNodeSpec `json:"objects",omitempty`
}

type TreeNodeSpec struct {
	ObjectID int64          `json:"objectid",omitempty`
	Name     string         `json:"name",omitempty`
	Objects  []TreeNodeSpec `json:"objects",omitempty`
}

// API struct holds all paths necessary to access Model Derivative API
type ModelDerivativeAPI struct {
	Client  *clientapi.Client
	APIPath string
}

// NewAPIWithCredentials returns a Model Derivative API client with default configurations
func NewAPIWithCredentials(ClientID string, ClientSecret string) ModelDerivativeAPI {
	auth := twolegged.NewAuth(ClientID, ClientSecret)
	return ModelDerivativeAPI{
		Client: clientapi.NewClient(auth),
	}
}

func (api ModelDerivativeAPI) path(paths ...string) []string {
	if api.APIPath == "" {
		return append([]string{DefaultModelDerivativePath}, paths...)
	}
	return append([]string{api.APIPath}, paths...)
}

// TranslateWithParams triggers translation job with settings specified in given TranslationParams
func (api ModelDerivativeAPI) TranslateWithParams(params TranslationParams) (result TranslationResult, err error) {
	byteParams, err := json.Marshal(params)
	if err != nil {
		return result, err
	}

	res, err := api.Client.DoRawRequest(
		context.Background(), http.MethodPost,
		scopes.DataRead|scopes.DataWrite,
		api.path("job"),
		nil, nil,
		clientapi.ContentTypeJSON,
		bytes.NewBuffer(byteParams),
	)
	if err != nil {
		return result, err
	}
	err = api.Client.ProcessRawError(res, &result)
	return result, err
}

// TranslateToSVF is a helper function that will use the TranslationSVFPreset for translating into svf a given ObjectID.
// It will also take care of converting objectID into Base64 (URL Safe) encoded URN.
func (api ModelDerivativeAPI) TranslateToSVF(objectID string) (result TranslationResult, err error) {
	params := TranslationSVFPreset
	params.Input.URN = base64.RawURLEncoding.EncodeToString([]byte(objectID))
	return api.TranslateWithParams(params)
}

func (api ModelDerivativeAPI) GetManifest(urn string) (result ManifestResult, err error) {
	res, err := api.Client.DoRawRequest(
		context.Background(), http.MethodGet,
		scopes.DataRead,
		api.path(urn, "manifest"),
		nil, nil,
		clientapi.ContentTypeJSON,
		nil,
	)
	if err != nil {
		return result, err
	}
	err = api.Client.ProcessRawError(res, &result)
	return result, err
}

func (api ModelDerivativeAPI) GetMetadata(urn string) (result MetadataResult, err error) {
	res, err := api.Client.DoRawRequest(
		context.Background(), http.MethodGet,
		scopes.DataRead,
		api.path(urn, "metadata"),
		nil, nil,
		clientapi.ContentTypeJSON,
		nil,
	)
	if err != nil {
		return result, err
	}
	err = api.Client.ProcessRawError(res, &result)
	return result, err
}

func (api ModelDerivativeAPI) GetObjectTree(urn string, viewID string) (status int, result TreeResult, err error) {

	res, err := api.Client.DoRawRequest(
		context.Background(), http.MethodGet,
		scopes.DataRead,
		api.path(urn, "metadata", viewID),
		[]clientapi.Filterer{filters.QueryParam{Key: "forceget", Value: "true"}},
		nil,
		clientapi.ContentTypeJSON,
		nil,
	)
	if err != nil {
		return 0, result, err
	}
	err = api.Client.ProcessRawError(res, &result)
	return res.StatusCode, result, err
}

func (api ModelDerivativeAPI) GetPropertiesStream(urn string, viewID string) (status int, result io.ReadCloser, err error) {
	res, err := api.Client.DoRawRequest(
		context.Background(), http.MethodGet,
		scopes.DataRead,
		api.path(urn, "metadata", viewID, "properties"),
		[]clientapi.Filterer{filters.QueryParam{Key: "forceget", Value: "true"}},
		nil,
		clientapi.ContentTypeJSON,
		nil,
	)
	if err != nil {
		return 0, nil, err
	}
	return res.StatusCode, res.Body, nil
}

func (api ModelDerivativeAPI) GetPropertiesObject(urn string, viewID string) (result PropertiesResult, err error) {

	status, stream, err := api.GetPropertiesStream(urn, viewID)
	if err != nil {
		return result, err
	}
	defer stream.Close()

	//using 200 as an error mask since it can be 2xx depending on state
	if (status & http.StatusOK) == 0 {
		content, _ := ioutil.ReadAll(stream)
		return result, clientapi.ErrorResult{StatusCode: status, Reason: string(content)}
	}
	decoder := json.NewDecoder(stream)
	err = decoder.Decode(&result)
	return result, err

}

func (api ModelDerivativeAPI) GetThumbnail(urn string) (reader io.ReadCloser, err error) {
	response, err := api.Client.DoRawRequest(
		context.Background(), http.MethodGet,
		scopes.DataRead,
		api.path(urn, "thumbnail"),
		nil, nil,
		clientapi.ContentTypeJSON,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		return nil, clientapi.ErrorResult{StatusCode: response.StatusCode, Reason: string(content)}
	}
	return response.Body, nil
}
