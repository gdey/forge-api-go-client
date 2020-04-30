package threelegged

import (
	"context"

	"github.com/gdey/forge-api-go-client/api"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

const (
	DefaultInformationalAPIPath = "userprofile/v1"
)

// UserProfile reflects the response received when query the profile of an authorizing end user in a 3-legged context
type UserProfile struct {
	UserID    string `json:"userId"`    // The backend user ID of the profile
	UserName  string `json:"userName"`  // The username chosen by the user
	EmailID   string `json:"emailId"`   // The user’s email address
	FirstName string `json:"firstName"` // The user’s first name
	LastName  string `json:"lastName"`  // The user’s last name
	// true if the user’s email address has been verified false if the user’s email address has not been verified
	EmailVerified bool `json:"emailVerified"`
	// true if the user has enabled two-factor authentication false if the user has not enabled two-factor authentication
	Var2FaEnabled bool `json:"2FaEnabled"`
	// A flat JSON object of attribute-value pairs in which the attributes specify available profile image sizes in the
	// format sizeX<pixels> (where <pixels> is an integer that represents both height and width in pixels of square
	// profile images) and the values are URLs for downloading the images via HTTP
	ProfileImages interface{} `json:"profileImages"`
}

// Information struct is holding the host and path used when making queries
// for profile of an authorizing end user in a 3-legged context
type Information struct {
	APIPath string
	AuthToken
}

func (info Information) Path(paths ...string) []string {
	if info.APIPath == "" {
		return append([]string{DefaultInformationalAPIPath}, paths...)
	}
	return append([]string{info.APIPath}, paths...)
}

//AboutMe is used to get the profile of an authorizing end user, given the token obtained via 3-legged OAuth flow
func (info Information) AboutMe() (profile UserProfile, err error) {

	client := api.NewClient(info.AuthToken)
	err = client.Get(
		context.Background(),
		scopes.UserProfileRead,
		info.Path("users/@me"),
		&profile,
		nil,
	)
	return profile, err
}
