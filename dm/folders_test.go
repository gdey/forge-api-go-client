package dm_test

import (
	"os"
	"testing"

	"github.com/outer-labs/forge-api-go-client/dm"
	"github.com/outer-labs/forge-api-go-client/env"
)

func TestFolderAPI_GetFolderDetails(t *testing.T) {

	// prepare the credentials
	clientID, clientSecret := env.GetClientSecretTest(t)

	folderAPI := dm.NewFolderAPIWithCredentials(clientID, clientSecret)

	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")
	testFolderKey := os.Getenv("BIM_360_TEST_ACCOUNT_FOLDERKEY")

	if testFolderKey == "" || testProjectKey == "" {
		t.Skip("Test environment vars not set")
	}

	t.Run("List all folders for a given project", func(t *testing.T) {
		_, err := folderAPI.GetFolderDetails(testProjectKey, testFolderKey)

		if err != nil {
			t.Fatalf("Failed to get project details: %s\n", err.Error())
		}
	})
}

func TestFolderAPI_GetContents(t *testing.T) {

	// prepare the credentials
	clientID, clientSecret := env.GetClientSecretTest(t)

	folderAPI := dm.NewFolderAPIWithCredentials(clientID, clientSecret)

	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")
	testFolderKey := os.Getenv("BIM_360_TEST_ACCOUNT_FOLDERKEY")

	t.Run("Get folder contents", func(t *testing.T) {
		_, err := folderAPI.GetFolderContents(testProjectKey, testFolderKey)

		if err != nil {
			t.Fatalf("Failed to get folder contents: %s\n", err.Error())
		}
	})

	t.Run("Get nonexistent folder contents", func(t *testing.T) {
		_, err := folderAPI.GetFolderContents(testProjectKey, testFolderKey+"30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting details for non-existing folder contents\n")
		}
	})
}
