package retoolsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client is the main struct for the SDK
// APIKey is the API key for the Retool instance (required)
// Endpoint is the URL of the Retool instance (required)
// Protocol is the protocol to use for the API requests (default: https).
type Client struct {
	APIKey     string
	Endpoint   string
	BaseURL    string
	HTTPClient *http.Client
}

// Response is the struct for the response from the Retool API
// By default, responses include up to 100 items.
// When there are more items, the has_more field in the response is set to true
// and the next field has a pagination token.
type Response[T any] struct {
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
	Data       T      `json:"data,omitempty"`
	TotalCount int    `json:"total_count,omitempty"`
	NextToken  string `json:"next_token,omitempty"`
	HasMore    bool   `json:"has_more,omitempty"`
}

// transportWithAPIKey is custom transport that adds the API key to
// the Authorization header for every request.
type transportWithAPIKey struct {
	APIKey    string
	Transport http.RoundTripper
}

// ClientOption defines the type for functional options.
type ClientOption func(*Client) error

// NewClient creates a new Retool client
// apiKey is the API key for the Retool instance (required)
// Endpoint is the URL of the Retool instance (required) (default: https)
// opts are the optional client options.
func NewClient(apiKey, endpoint string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" || endpoint == "" {
		return nil, errors.New("API key and Endpoint are required")
	}

	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	customTransport := &transportWithAPIKey{
		APIKey:    apiKey,
		Transport: http.DefaultTransport,
	}

	c := &Client{
		APIKey:   apiKey,
		Endpoint: endpoint,
		BaseURL:  endpoint + "/api/v2",
		HTTPClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: customTransport,
		},
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("applying client option: %w", err)
		}
	}

	return c, nil
}

// Do makes an HTTP request to the Retool API.
func (c *Client) Do(method, url string, body interface{}) (*http.Response, error) {
	var requestBody []byte
	var err error

	if body != nil {
		requestBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshalling request: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	return resp, nil
}

// RoundTrip adds the API key to the Authorization header for every request
// and sets the Content-Type header to application/json.
func (t *transportWithAPIKey) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.APIKey))

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return t.Transport.RoundTrip(req)
}

// WithTimeout allows setting a custom timeout in seconds for the HTTP client.
// Have to pass in the time.Duration value with time.Type (e.g. 10*time.Second).
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		if timeout > 0 {
			c.HTTPClient.Timeout = timeout
		} else {
			return errors.New("timeout must be greater than 0")
		}
		return nil
	}
}
