package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Request is the data structure used to request a token from the Airthings API.
type Request struct {
	GrantType    string   `json:"grant_type"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scope        []string `json:"scope"`
}

// IsValid checks if the request is valid.
func (r *Request) IsValid() bool {
	return r.GrantType != "" && r.ClientID != "" && r.ClientSecret != "" && len(r.Scope) > 0
}

// Token is the data structure used to represent a token from the Airthings API.
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Fetcher struct {
	Client HTTPClient
	URL    string
}

func (f *Fetcher) GetToken(req *Request) (*Token, error) {
	if req == nil || !req.IsValid() {
		return nil, fmt.Errorf("invalid request")
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request JSON: %w", err)
	}

	r, err := http.NewRequest(http.MethodGet, f.URL, bytes.NewBuffer(jsonData))
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
