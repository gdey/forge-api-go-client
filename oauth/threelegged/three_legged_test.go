package threelegged_test

import (
	"testing"

	"github.com/gdey/forge-api-go-client/env"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
	"github.com/gdey/forge-api-go-client/oauth/threelegged"
)

func TestAuth_Authorize(t *testing.T) {

	//prepare the credentials
	clientID, clientSecret := env.GetClientSecretTest(t)

	client := threelegged.NewAuth(clientID,
		clientSecret,
		"http://localhost:3009/callback",
		scopes.DataRead|scopes.DataWrite,
	)

	authLink, err := client.Authorize("something that will be passed back")

	if err != nil {
		t.Errorf("Could not create the authorization link, got: %s", err.Error())
	}

	if len(authLink) == 0 {
		t.Errorf("The authorization link is empty")
	}

}
