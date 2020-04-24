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
	"github.com/gdey/forge-api-go-client/oauth"
)

const (
	DefaultRecapAPIPath = "/photo-to-3d/v1"
)

// API struct holds all paths necessary to access ReCap API
type API struct {
	oauth.ForgeAuthenticator
	APIPath string
}

// NewAPIWithCredentials returns a ReCap API client with default configurations
func NewAPIWithCredentials(ClientID string, ClientSecret string) API {
	return API{
		ForgeAuthenticator: oauth.NewTwoLeggedClient(ClientID, ClientSecret),
	}
}

func (api API) Path() string {
	if api.APIPath != "" {
		return api.ForgeAuthenticator.HostPath(api.APIPath)
	}
	return api.ForgeAuthenticator.HostPath(DefaultRecapAPIPath)
}

// CreatePhotoScene prepares a scene with a given name, expected output formats and sceneType
// 	name - should not be empty
// 	formats - should be of type rcm, rcs, obj, ortho or report
// 	sceneType - should be either "aerial" or "object"
func (api API) CreatePhotoScene(name string, formats []string, sceneType string) (scene PhotoScene, err error) {
	bearer, err := api.GetTokenWithScope(oauth.ScopeDataWrite)
	if err != nil {
		return
	}
	return createPhotoScene(api.Path(), name, formats, sceneType, bearer.AccessToken)
}

// AddFileToSceneUsingLink can be used when the needed images are already available remotely
// and can be uploaded just by providing the remote link
func (api API) AddFileToSceneUsingLink(sceneID string, link string) (uploads FileUploadingReply, err error) {
	bearer, err := api.GetTokenWithScope(oauth.ScopeDataWrite)
	if err != nil {
		return
	}
	return addFileToSceneUsingLink(api.Path(), sceneID, link, bearer.AccessToken)
}

// AddFileToSceneUsingData can be used when the image is already available as a byte slice,
// be it read from a local file or as a result/body of a POST request
func (api API) AddFileToSceneUsingData(sceneID string, data []byte) (uploads FileUploadingReply, err error) {
	bearer, err := api.GetTokenWithScope(oauth.ScopeDataWrite)
	if err != nil {
		return
	}
	return addFileToSceneUsingFileData(api.Path(), sceneID, data, bearer.AccessToken)
}

// StartSceneProcessing will trigger the processing of a specified scene that can be canceled any time
func (api API) StartSceneProcessing(sceneID string) (result SceneStartProcessingReply, err error) {
	bearer, err := api.GetTokenWithScope(oauth.ScopeDataWrite)
	if err != nil {
		return
	}
	return startSceneProcessing(api.Path(), sceneID, bearer.AccessToken)
}

// GetSceneProgress polls the scene processing status and progress
//	Note: instead of polling, consider using the callback parameter that can be specified upon scene creation
func (api API) GetSceneProgress(sceneID string) (progress SceneProgressReply, err error) {
	bearer, err := api.GetTokenWithScope(oauth.ScopeDataRead)
	if err != nil {
		return
	}
	return getSceneProgress(api.Path(), sceneID, bearer.AccessToken)
}

// GetSceneResults requests result in a specified format
//	Note: The link specified in SceneResultReplies will be available for the time specified in reply,
//	even if the scene is deleted
func (api API) GetSceneResults(sceneID string, format string) (result SceneResultReply, err error) {
	bearer, err := api.GetTokenWithScope(oauth.ScopeDataRead)
	if err != nil {
		return
	}
	return getSceneResult(api.Path(), sceneID, bearer.AccessToken, format)
}

// CancelSceneProcessing stops the scene processing, without affecting the already uploaded resources
func (api API) CancelSceneProcessing(sceneID string) (ID string, err error) {
	bearer, err := api.GetTokenWithScope(oauth.ScopeDataWrite)
	if err != nil {
		return sceneID, err
	}
	_, err = cancelSceneProcessing(api.Path(), sceneID, bearer.AccessToken)

	return sceneID, err
}

// DeleteScene removes all the resources associated with given scene.
func (api API) DeleteScene(sceneID string) (ID string, err error) {
	bearer, err := api.GetTokenWithScope(oauth.ScopeDataWrite)
	if err != nil {
		return
	}
	_, err = deleteScene(api.Path(), sceneID, bearer.AccessToken)
	return sceneID, nil
}
