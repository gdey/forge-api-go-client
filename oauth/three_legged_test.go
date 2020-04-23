package oauth_test

import (
	"testing"

	"github.com/outer-labs/forge-api-go-client/env"
	"github.com/outer-labs/forge-api-go-client/oauth"
)

func TestThreeLeggedAuth_Authorize(t *testing.T) {

	//prepare the credentials
	clientID, clientSecret := env.GetClientSecretTest(t)

	client := oauth.NewThreeLeggedClient(clientID,
		clientSecret,
		"http://localhost:3009/callback")

	authLink, err := client.Authorize("data:read data:write", "something that will be passed back")

	if err != nil {
		t.Errorf("Could not create the authorization link, got: %s", err.Error())
	}

	if len(authLink) == 0 {
		t.Errorf("The authorization link is empty")
	}

}
