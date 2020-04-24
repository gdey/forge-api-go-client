package oauth_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gdey/forge-api-go-client/env"
	"github.com/gdey/forge-api-go-client/oauth"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

func TestAuthenticate(t *testing.T) {

	clientID, clientSecret := env.GetClientSecretTest(t)

	t.Run("Valid Forge Secrets", func(t *testing.T) {
		authenticator := oauth.NewTwoLeggedClient(clientID, clientSecret)

		bearer, err := authenticator.Authenticate(scopes.DataRead)

		if err != nil {
			t.Error(err.Error())
		}

		if len(bearer.AccessToken) == 0 {
			t.Errorf("Wrong bearer content: %v", bearer)
		}
	})

	t.Run("Invalid Forge Secrets", func(t *testing.T) {
		authenticator := oauth.NewTwoLeggedClient("", clientSecret)

		bearer, err := authenticator.Authenticate(scopes.DataRead)

		if err == nil {
			t.Errorf("Expected to fail due to wrong credentials, but got %v", bearer)
		}

		if len(bearer.AccessToken) != 0 {
			t.Errorf("expected to not receive a token, but received: %s", bearer.AccessToken)
		}
	})

	t.Run("Invalid scope", func(t *testing.T) {
		authenticator := oauth.NewTwoLeggedClient(clientID, clientSecret)

		// Get a bad scope. if ScopeAccountWrite is not the last scope, then this will fail
		var badScope scopes.Scope
		bearer, err := authenticator.Authenticate(badScope)

		if err == nil {
			t.Errorf("Expected to fail due to wrong scope, but got %v\n", bearer)
		}

		if len(bearer.AccessToken) != 0 {
			t.Errorf("expected to not receive a token, but received: %s", bearer.AccessToken)
		}
	})

	t.Run("Invalid or unreachable host", func(t *testing.T) {
		authenticator := oauth.NewTwoLeggedClient(clientID, clientSecret)
		authenticator.Host = "http://localhost"

		bearer, err := authenticator.Authenticate(scopes.DataRead)

		if err == nil {
			t.Errorf("Expected to fail due to wrong host, but got %v\n", bearer)
		}

		if len(bearer.AccessToken) != 0 {
			t.Errorf("expected to not receive a token, but received: %s", bearer.AccessToken)
		}
	})
}

func ExampleTwoLeggedAuth_Authenticate() {

	// acquire Forge secrets from environment
	clientID, clientSecret := os.Getenv("FORGE_CLIENT_ID"), os.Getenv("FORGE_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		log.Fatalf("Could not get from env the Forge secrets")
	}

	// create oauth client
	authenticator := oauth.NewTwoLeggedClient(clientID, clientSecret)

	// request a token with needed scopes, separated by spaces
	bearer, err := authenticator.Authenticate(scopes.DataRead | scopes.DataWrite)

	if err != nil || len(bearer.AccessToken) == 0 {
		log.Fatalf("Could not get from env the Forge secrets")
	}

	// at this point, the bearer should contain the needed data. Check Bearer struct for more info
	fmt.Printf("Bearer now contains:\n"+
		"AccessToken: %s\n"+
		"TokenType: %s\n"+
		"Expires in: %d\n",
		bearer.AccessToken,
		bearer.TokenType,
		bearer.ExpiresIn)

}
