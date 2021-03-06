package sethealth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

const host = "https://api.set.health"

type GetTokenOptions struct {
	UserID    string
	ExpiresIn time.Duration
	TestMode  bool
}

// Client exposes the public api for sethealth
type Client struct {
	key    string
	secret string
	client *http.Client
}

var ErrCredentials = errors.New("Invalid credentials")

// GetTokenResponse contains the result values of calling the GetToken() API
type GetTokenResponse struct {
	Token string `json:"token"`
}

// New creates a new client for the server sethealth API
// It will automatically get the Sethealth credentials from the
// "SETHEALTH_KEY" and "SETHEALTH_SECRET" environment variables.
func New() *Client {
	return NewWithCredentials(
		os.Getenv("SETHEALTH_KEY"),
		os.Getenv("SETHEALTH_SECRET"),
	)
}

// NewWithCredentials creates a new client for the server sethealth API
// It requires a key and secret in order to perform any request
func NewWithCredentials(key, secret string) *Client {
	return &Client{
		key:    key,
		secret: secret,
		client: &http.Client{},
	}
}

// GetToken returns a new short-living token to be used by client side.
func (c *Client) GetToken() (GetTokenResponse, error) {
	return c.GetTokenWithOptions(GetTokenOptions{})
}

// GetTokenWithOptions returns a new short-living token to be used by client side with options.
func (c *Client) GetTokenWithOptions(opts GetTokenOptions) (GetTokenResponse, error) {
	var token GetTokenResponse
	data := map[string]interface{}{
		"id":         c.key,
		"secret":     c.secret,
		"test-mode":  opts.TestMode,
		"expires-in": opts.ExpiresIn,
		"user-id":    opts.UserID,
	}
	jsonBytes, _ := json.Marshal(data)
	body := bytes.NewBuffer(jsonBytes)

	// Create request
	req, err := http.NewRequest("POST", host+"/token", body)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// Fetch Request
	resp, err := c.client.Do(req)
	if err != nil {
		return token, err
	}
	if resp.StatusCode != 200 {
		return token, ErrCredentials
	}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return token, err
	}
	return token, nil
}
