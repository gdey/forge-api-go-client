package dm_test

import (
	"os"
	"testing"

	"github.com/outer-labs/forge-api-go-client/dm"
	"github.com/outer-labs/forge-api-go-client/env"
)

func TestProjectAPI_GetProjects(t *testing.T) {

	// prepare the credentials
	clientID, clientSecret := env.GetClientSecretTest(t)

	hubAPI := dm.NewHubAPIWithCredentials(clientID, clientSecret)

	testHubKey := os.Getenv("BIM_360_TEST_ACCOUNT_HUBKEY")

	t.Run("List all projects under a given hub", func(t *testing.T) {
		_, err := hubAPI.ListProjects(testHubKey)

		if err != nil {
			t.Fatalf("Failed to get project details: %s\n", err.Error())
		}
	})

	t.Run("List all projects under non-existent hub (should fail)", func(t *testing.T) {
		_, err := hubAPI.ListProjects(testHubKey + "30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting projects for non-existing hub\n")
		}
	})
}

func TestProjectAPI_GetProjectDetails(t *testing.T) {

	// prepare the credentials
	clientID, clientSecret := env.GetClientSecretTest(t)

	hubAPI := dm.NewHubAPIWithCredentials(clientID, clientSecret)

	testHubKey := os.Getenv("BIM_360_TEST_ACCOUNT_HUBKEY")
	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")

	t.Run("List all projects under a given hub", func(t *testing.T) {
		_, err := hubAPI.GetProjectDetails(testHubKey, testProjectKey)

		if err != nil {
			t.Fatalf("Failed to get project details: %s\n", err.Error())
		}
	})

	t.Run("List all projects under non-existent hub (should fail)", func(t *testing.T) {
		_, err := hubAPI.GetProjectDetails(testHubKey, testProjectKey+"30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting projects for non-existing hub\n")
		}
	})
}

func TestProjectAPI_GetTopFolders(t *testing.T) {

	// prepare the credentials
	clientID, clientSecret := env.GetClientSecretTest(t)

	hubAPI := dm.NewHubAPIWithCredentials(clientID, clientSecret)

	testHubKey := os.Getenv("BIM_360_TEST_ACCOUNT_HUBKEY")
	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")

	t.Run("List all projects under a given hub", func(t *testing.T) {
		_, err := hubAPI.GetTopFolders(testHubKey, testProjectKey)

		if err != nil {
			t.Fatalf("Failed to get project details: %s\n", err.Error())
		}
	})

	t.Run("List all projects under non-existent hub (should fail)", func(t *testing.T) {
		_, err := hubAPI.GetTopFolders(testHubKey, testProjectKey+"30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting projects for non-existing hub\n")
		}
	})
}
