package threelegged_test

import (
	"fmt"

	"github.com/gdey/forge-api-go-client/oauth"

	"github.com/gdey/forge-api-go-client/oauth/threelegged"
)

//TODO: enable it after having set up a pipeline for auto-creating a 3-legged oauth token
//func TestInformation_AboutMe(t *testing.T) {
//
//	info := oauth.NewInformationQuerier()
//
//	aThreeLeggedToken := os.Getenv("THREE_LEGGED_TOKEN")
//
//	profile, err := info.AboutMe(aThreeLeggedToken)
//
//	if err != nil {
//		t.Error(err.Error())
//		return
//	}
//
//	t.Logf("Received profile:\n"+
//		"UserId: %s\n"+
//		"UserName: %s\n"+
//		"EmailId: %s\n"+
//		"FirstName: %s\n"+
//		"LastName: %s\n"+
//		"EmailVerified: %t\n"+
//		"Var2FaEnabled: %t\n"+
//		"ProfileImages: %v",
//		profile.UserID,
//		profile.UserName,
//		profile.EmailID,
//		profile.FirstName,
//		profile.LastName,
//		profile.EmailVerified,
//		profile.Var2FaEnabled,
//		profile.ProfileImages)
//}

func ExampleInformation_AboutMe() {

	aThreeLeggedToken := "put a valid 3-legged token here"
	Auth := threelegged.AuthToken{
		Token: threelegged.NewRefreshableToken(&oauth.Bearer{
			TokenType:   "Bearer",
			ExpiresIn:   60,
			AccessToken: aThreeLeggedToken,
		}),
	}

	info := threelegged.Information{
		AuthToken: Auth,
	}

	profile, err := info.AboutMe()

	if err != nil {
		fmt.Printf("[ERROR] Could not retrieve profile, got %s\n", err.Error())
		return
	}

	fmt.Printf("Received profile:\n"+
		"UserId: %s\n"+
		"UserName: %s\n"+
		"EmailId: %s\n"+
		"FirstName: %s\n"+
		"LastName: %s\n"+
		"EmailVerified: %t\n"+
		"Var2FaEnabled: %t\n"+
		"ProfileImages: %v",
		profile.UserID,
		profile.UserName,
		profile.EmailID,
		profile.FirstName,
		profile.LastName,
		profile.EmailVerified,
		profile.Var2FaEnabled,
		profile.ProfileImages,
	)
}
