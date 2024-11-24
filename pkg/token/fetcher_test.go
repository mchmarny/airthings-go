package token

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// MockHTTPClient simulates an HTTP client for testing.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestTokenFetcher_GetToken(t *testing.T) {
	mockClient := &MockHTTPClient{}

	r := &Request{
		GrantType:    "client_credentials",
		ClientID:     "test_id",
		ClientSecret: "test_secret",
		Scope:        []string{"read:device:current_values"},
	}

	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data := string(b)

	// Simulate a successful token response
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Errorf("expected POST method, got %s", req.Method)
		}
		if req.URL.String() != "https://accounts-api.airthings.com/v1/token" {
			t.Errorf("unexpected URL: %s", req.URL.String())
		}
		body, _ := io.ReadAll(req.Body)
		expectedBody := data
		if string(body) != expectedBody {
			t.Errorf("unexpected request body: %s", string(body))
		}

		responseBody := `{
			"access_token": "test_token",
			"expires_in": 3600,
			"token_type": "Bearer"
		}`

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		}, nil
	}

	f := &Fetcher{
		Client: mockClient,
		URL:    "https://accounts-api.airthings.com/v1/token",
	}

	token, err := f.GetToken(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token.AccessToken != "test_token" {
		t.Errorf("unexpected access token: %s", token.AccessToken)
	}
	if token.ExpiresIn != 3600 {
		t.Errorf("unexpected expires_in: %d", token.ExpiresIn)
	}
	if token.TokenType != "Bearer" {
		t.Errorf("unexpected token_type: %s", token.TokenType)
	}
}

func TestTokenFetcher_GetToken_ErrorCases(t *testing.T) {
	mockClient := &MockHTTPClient{}

	// Simulate an API error response
	mockClient.DoFunc = func(_ *http.Request) (*http.Response, error) {
		responseBody := `{"error": "invalid_client"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		}, nil
	}

	f := &Fetcher{
		Client: mockClient,
		URL:    "https://accounts-api.airthings.com/v1/token",
	}

	_, err := f.GetToken(&Request{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "invalid request" {
		t.Errorf("unexpected error message: %v", err)
	}
}
