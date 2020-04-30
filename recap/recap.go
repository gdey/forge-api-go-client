// Package recap contains the Go wrappers for calls to Forge Reality Capture API
// https://developer.autodesk.com/api/reality-capture-cover-page/
//
// 	The workflow is simple:
// 		- create a photoScene
//		- upload images to photoScene
//		- start photoScene processing
//		- get the result
package recap

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/gdey/forge-api-go-client/api/filters"

	clientapi "github.com/gdey/forge-api-go-client/api"
	"github.com/gdey/forge-api-go-client/oauth"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
	"github.com/gdey/forge-api-go-client/oauth/twolegged"
)

const (
	DefaultRecapAPIPath = "/photo-to-3d/v1"
)

// API struct holds all paths necessary to access ReCap API
type API struct {
	Client  *clientapi.Client
	APIPath string

	oauth.ForgeAuthenticator
}

// NewAPIWithCredentials returns a ReCap API client with default configurations
func NewAPIWithCredentials(ClientID string, ClientSecret string) API {
	auth := twolegged.NewAuth(ClientID, ClientSecret)
	return API{
		ForgeAuthenticator: auth,
		Client:             clientapi.NewClient(auth),
	}
}

func (api API) Path(paths ...string) []string {
	if api.APIPath == "" {
		return append([]string{DefaultRecapAPIPath}, paths...)
	}
	return append([]string{api.APIPath}, paths...)
}

// CreatePhotoScene prepares a scene with a given name, expected output formats and sceneType
// 	name - should not be empty
// 	formats - should be of type rcm, rcs, obj, ortho or report
// 	sceneType - should be either "aerial" or "object"
func (api API) CreatePhotoScene(name string, formats []string, sceneType string) (scene PhotoScene, err error) {
	// TODO(gdey): sceneType should be a custom type
	if sceneType != "object" && sceneType != "aerial" {
		err = errors.New("the scene type is not supported. Expecting 'object' or 'aerial', got " + sceneType)
		return
	}
	body := url.Values{
		"scenename": []string{name},
		"format":    []string{strings.Join(formats, ",")},
		"scenetype": []string{sceneType},
	}
	response, err := api.Client.DoRawRequest(
		context.Background(), http.MethodPost,
		scopes.DataWrite,
		api.Path("photoscene"),
		nil,
		nil,
		clientapi.ContentTypeJSON,
		bytes.NewBufferString(body.Encode()),
	)
	if err != nil {
		return scene, err
	}
	defer response.Body.Close()
	var result SceneCreationReply
	if err = api.Client.ProcessRawError(response, &result); err != nil {
		return scene, err
	}

	// This check is necessary, as there are cases when server returns status OK, but contains an error message
	if result.Error != nil {
		return scene, result.Error
	}
	return result.PhotoScene, nil

}

// AddFileToSceneUsingLink can be used when the needed images are already available remotely
// and can be uploaded just by providing the remote link
func (api API) AddFileToSceneUsingLink(sceneID string, link string) (uploads FileUploadingReply, err error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writer.WriteField("photosceneid", sceneID)
	writer.WriteField("type", "image")
	writer.WriteField("file[0]", link)
	writer.Close()
	response, err := api.Client.DoRawRequest(
		context.Background(), http.MethodPost,
		scopes.DataWrite,
		api.Path("file"),
		nil,
		nil,
		clientapi.ContentTypeJSON,
		&body,
	)
	if err != nil {
		return uploads, err
	}
	defer response.Body.Close()
	if err = api.Client.ProcessRawError(response, &uploads); err != nil {
		err = fmt.Errorf("ProcessRawError: %w", err)
		return uploads, err
	}
	// This check is necessary, as there are cases when server returns status OK, but contains an error message
	if uploads.Error != nil {
		return uploads, uploads.Error
	}
	return uploads, nil

}

// AddFileToSceneUsingData can be used when the image is already available as a byte slice,
// be it read from a local file or as a result/body of a POST request
func (api API) AddFileToSceneUsingData(sceneID string, data []byte) (uploads FileUploadingReply, err error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writer.WriteField("photosceneid", sceneID)
	writer.WriteField("type", "image")
	formFile, err := writer.CreateFormFile("file[0]", fmt.Sprintf("data%d", rand.Int()))
	if err != nil {
		return uploads, err
	}
	formFile.Write(data)
	writer.Close()

	response, err := api.Client.DoRawRequest(
		context.Background(), http.MethodPost,
		scopes.DataWrite,
		api.Path("file"),
		nil,
		nil,
		clientapi.ContentTypeJSON,
		&body,
	)
	if err != nil {
		return uploads, err
	}
	defer response.Body.Close()
	if err = api.Client.ProcessRawError(response, &uploads); err != nil {
		return uploads, err
	}
	// This check is necessary, as there are cases when server returns status OK, but contains an error message
	if uploads.Error != nil {
		return uploads, uploads.Error
	}
	return uploads, nil
}

// StartSceneProcessing will trigger the processing of a specified scene that can be canceled any time
func (api API) StartSceneProcessing(sceneID string) (result SceneStartProcessingReply, err error) {
	response, err := api.Client.DoRawRequest(
		context.Background(), http.MethodPost,
		scopes.DataWrite,
		api.Path("photoscene", sceneID),
		nil, nil, clientapi.ContentTypeJSON, nil,
	)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()
	if err = api.Client.ProcessRawError(response, &result); err != nil {
		return result, err
	}
	// This check is necessary, as there are cases when server returns status OK, but contains an error message
	if result.Error != nil {
		return result, result.Error
	}
	return result, nil
}

// GetSceneProgress polls the scene processing status and progress
//	Note: instead of polling, consider using the callback parameter that can be specified upon scene creation
func (api API) GetSceneProgress(sceneID string) (progress SceneProgressReply, err error) {
	response, err := api.Client.DoRawRequest(
		context.Background(), http.MethodGet,
		scopes.DataRead,
		api.Path("photoscene", sceneID, "progress"),
		nil, nil, clientapi.ContentTypeJSON, nil,
	)
	if err != nil {
		return progress, err
	}
	defer response.Body.Close()
	if err = api.Client.ProcessRawError(response, &progress); err != nil {
		return progress, err
	}
	// This check is necessary, as there are cases when server returns status OK, but contains an error message
	if progress.Error != nil {
		return progress, progress.Error
	}
	return progress, nil
}

// GetSceneResults requests result in a specified format
//	Note: The link specified in SceneResultReplies will be available for the time specified in reply,
//	even if the scene is deleted
func (api API) GetSceneResults(sceneID string, format string) (result SceneResultReply, err error) {
	response, err := api.Client.DoRawRequest(
		context.Background(), http.MethodGet,
		scopes.DataRead,
		api.Path("photoscene", sceneID),
		[]clientapi.Filterer{filters.QueryParam{Key: "format", Value: format}},
		nil, clientapi.ContentTypeJSON, nil,
	)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()
	if err = api.Client.ProcessRawError(response, &result); err != nil {
		return result, err
	}
	// This check is necessary, as there are cases when server returns status OK, but contains an error message
	if result.Error != nil {
		return result, result.Error
	}
	return result, nil
}

// CancelSceneProcessing stops the scene processing, without affecting the already uploaded resources
func (api API) CancelSceneProcessing(sceneID string) (ID string, err error) {
	var result SceneCancelReply
	response, err := api.Client.DoRawRequest(
		context.Background(), http.MethodGet,
		scopes.DataWrite,
		api.Path("photoscene", sceneID, "cancel"),
		nil,
		nil, clientapi.ContentTypeFormEncoded, nil,
	)
	if err != nil {
		return sceneID, err
	}
	defer response.Body.Close()
	if err = api.Client.ProcessRawError(response, &result); err != nil {
		return sceneID, err
	}
	// This check is necessary, as there are cases when server returns status OK, but contains an error message
	if result.Error != nil {
		return sceneID, result.Error
	}
	return sceneID, nil
}

// DeleteScene removes all the resources associated with given scene.
func (api API) DeleteScene(sceneID string) (ID string, err error) {
	return sceneID, nil
	var result SceneDeletionReply
	response, err := api.Client.DoRawRequest(
		context.Background(), http.MethodDelete,
		scopes.DataWrite,
		api.Path("photoscene", sceneID),
		nil,
		nil, clientapi.ContentTypeFormEncoded, nil,
	)
	if err != nil {
		return sceneID, err
	}
	defer response.Body.Close()
	if err = api.Client.ProcessRawError(response, &result); err != nil {
		return sceneID, err
	}
	// This check is necessary, as there are cases when server returns status OK, but contains an error message
	if result.Error != nil {
		return sceneID, result.Error
	}
	return sceneID, nil
}
