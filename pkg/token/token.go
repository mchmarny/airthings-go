package token

import "os"

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

// GetToken retrieves a token from the Airthings API using the client ID and secret environment variables.
func GetToken() (*Token, error) {
	return NewFetcher().GetToken(&Request{
		GrantType:    tokenGrantType,
		ClientID:     os.Getenv(ClientIDEnvVar),
		ClientSecret: os.Getenv(ClientSecretEnvVar),
		Scope:        []string{"read:device:current_values"},
	})
}
