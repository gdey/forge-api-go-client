package env

import (
	"os"
	"testing"
)

// GetClientSecretTest will retrive the ClientID and ClientSecret from the env, or
// call t.Skip if not found
func GetClientSecretTest(t *testing.T) (id, secret string) {
	// prepare the credentials
	id, secret = os.Getenv("FORGE_CLIENT_ID"), os.Getenv("FORGE_CLIENT_SECRET")
	if id == "" {
		t.Skip("ClientID not set")
	}
	if secret == "" {
		t.Skip("ClientSecret not set")
	}
	return id, secret
}

// GetClientSecret will retrive the ClientID and ClientSecret from the env
func GetClientSecret() (id, secret string) {
	// prepare the credentials
	return os.Getenv("FORGE_CLIENT_ID"), os.Getenv("FORGE_CLIENT_SECRET")
}

// GetTest will retrive the named env bar or Skip if
// the environmental variable is empty.
func GetTest(t *testing.T, name string) (value string) {
	value = os.Getenv(name)
	if value == "" {
		t.Skipf("%v not set", name)
	}
	return value
}
