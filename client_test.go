package retoolsdk_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	retool "github.com/thoughtgears/retoolsdk"

	"github.com/stretchr/testify/assert"
)

type MockTransport struct {
	Response *http.Response
	Err      error
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestNewClient_Success(t *testing.T) {
	client, err := retool.NewClient("test-api-key", "example.com")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, 10*time.Second, client.HTTPClient.Timeout)
}

func TestNewClient_MissingAPIKey(t *testing.T) {
	client, err := retool.NewClient("", "example.com")
	assert.Nil(t, client)
	assert.EqualError(t, err, "API key and Endpoint are required")
}

func TestNewClient_MissingEndpoint(t *testing.T) {
	client, err := retool.NewClient("test-api-key", "")
	assert.Nil(t, client)
	assert.EqualError(t, err, "API key and Endpoint are required")
}

func TestNewClient_CustomTimeout(t *testing.T) {
	client, err := retool.NewClient("test-api-key", "example.com", retool.WithTimeout(30*time.Second))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, 30*time.Second, client.HTTPClient.Timeout)
}

func TestNewClient_InvalidTimeout(t *testing.T) {
	client, err := retool.NewClient("test-api-key", "example.com", retool.WithTimeout(0))
	assert.Nil(t, client)
	assert.EqualError(t, err, "applying client option: timeout must be greater than 0")
}

func TestNewClient_RoundTrip(t *testing.T) {
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("OK")), // Simplified body
	}

	mockTransport := &MockTransport{
		Response: mockResponse,
	}

	client := &retool.Client{
		APIKey:   "test-api-key",
		Endpoint: "example.com",
		HTTPClient: &http.Client{
			Transport: mockTransport,
		},
	}

	req, _ := http.NewRequest("GET", "https://example.com/test", nil)
	resp, err := client.HTTPClient.Do(req)
	body, _ := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "OK", string(body))
}

func TestNewClient_AuthHeader(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client, err := retool.NewClient("test-api-key", mockServer.URL)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", mockServer.URL, nil)
	_, err = client.HTTPClient.Do(req)
	assert.NoError(t, err)
}
