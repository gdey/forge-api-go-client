package dm

import (
	"os"
	"testing"

	"github.com/gdey/forge-api-go-client/env"
)

func TestFolderAPI_GetItemDetails(t *testing.T) {

	// prepare the credentials
	clientID, clientSecret := env.GetClientSecretTest(t)

	folderAPI := NewFolderAPIWithCredentials(clientID, clientSecret)

	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")
	testItemKey := os.Getenv("BIM_360_TEST_ACCOUNT_ITEMKEY")

	t.Run("List item details", func(t *testing.T) {
		_, err := folderAPI.GetItemDetails(testProjectKey, testItemKey)

		if err != nil {
			t.Fatalf("Failed to get item details: %s\n", err.Error())
		}
	})
}

func TestFolderAPI_GetItemTip(t *testing.T) {

	// prepare the credentials
	clientID, clientSecret := env.GetClientSecretTest(t)

	folderAPI := NewFolderAPIWithCredentials(clientID, clientSecret)

	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")
	testItemKey := os.Getenv("BIM_360_TEST_ACCOUNT_ITEMKEY")

	t.Run("List item details", func(t *testing.T) {
		_, err := folderAPI.GetItemTip(testProjectKey, testItemKey)

		if err != nil {
			t.Fatalf("Failed to get item details: %s\n", err.Error())
		}
	})
}

func TestFolderAPI_GetItemVersions(t *testing.T) {

	// prepare the credentials
	clientID, clientSecret := env.GetClientSecretTest(t)

	folderAPI := NewFolderAPIWithCredentials(clientID, clientSecret)

	testProjectKey := os.Getenv("BIM_360_TEST_ACCOUNT_PROJECTKEY")
	testItemKey := os.Getenv("BIM_360_TEST_ACCOUNT_ITEMKEY")

	t.Run("List item details", func(t *testing.T) {
		_, err := folderAPI.GetItemVersions(testProjectKey, testItemKey)

		if err != nil {
			t.Fatalf("Failed to get item details: %s\n", err.Error())
		}
	})
}
