package client

import (
	"context"
	"fmt"
	"net/http"

	at "github.com/mchmarny/airthings-go/pkg/airthings"
	"github.com/mchmarny/airthings-go/pkg/token"
)

const (
	serverURL = "https://consumer-api.airthings.com/v1"
)

func NewClient() *Client {
	return &Client{}
}

type Client struct {
}

func (c *Client) GetDevices() (*http.Response, error) {
	ctx := context.Background()

	tokenFn := func(_ context.Context, req *http.Request) error {
		t, err := token.GetToken()
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}
		req.Header.Add("Token", fmt.Sprintf("Bearer %s", t.AccessToken))
		return nil
	}

	ac, err := at.NewClient(serverURL, at.WithRequestEditorFn(tokenFn))
	if err != nil {
		return nil, fmt.Errorf("failed to create airthings client: %w", err)
	}

	r, err := ac.GetHealth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get health: %w", err)
	}

	return r, nil
}
