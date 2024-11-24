package token

import (
	"fmt"
	"net/http"
	"os"
)

const (
	ClientIDEnvVar     = "AIRTHINGS_CLIENT_ID"
	ClientSecretEnvVar = "AIRTHINGS_CLIENT_SECRET"

	apiURL         = "https://accounts-api.airthings.com/v1/token"
	tokenGrantType = "client_credentials"

	clientTimeoutDefault = 10
)

// GetToken retrieves a token from the Airthings API using the client ID and secret environment variables.
func GetToken() (*Token, error) {
	return GetTokenWithRequest(&Request{
		GrantType:    tokenGrantType,
		ClientID:     os.Getenv(ClientIDEnvVar),
		ClientSecret: os.Getenv(ClientSecretEnvVar),
		Scope:        []string{"read:device:current_values"},
	})
}

// GetTokenWithValues retrieves a token from the Airthings API using the provided client ID and secret.
func GetTokenWithRequest(req *Request) (*Token, error) {
	if req == nil || !req.IsValid() {
		return nil, fmt.Errorf("invalid request")
	}

	client := &http.Client{
		Timeout: clientTimeoutDefault,
	}

	f := &Fetcher{
		Client: client,
		URL:    apiURL,
	}

	return f.GetToken(req)
}
