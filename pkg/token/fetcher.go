package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	ClientIDEnvVar     = "AIRTHINGS_CLIENT_ID"
	ClientSecretEnvVar = "AIRTHINGS_CLIENT_SECRET"

	apiURL         = "https://accounts-api.airthings.com/v1/token"
	tokenGrantType = "client_credentials"

	clientTimeoutDefault = 60
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewFetcher creates a new Fetcher with a default HTTP client and API URL.
func NewFetcher() *Fetcher {
	return &Fetcher{
		Client: &http.Client{
			Timeout: clientTimeoutDefault * time.Second,
		},
	}
}

type Fetcher struct {
	Client HTTPClient
}

func (f *Fetcher) GetToken(req *Request) (*Token, error) {
	if req == nil || !req.IsValid() {
		return nil, fmt.Errorf("invalid request")
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request JSON: %w", err)
	}

	r, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := f.Client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, body)
	}

	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &token, nil
}
