package retoolsdk_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/thoughtgears/retoolsdk"
	"io"
	"net/http"
	"testing"
)

func TestGetSpace_Success(t *testing.T) {
	response := `{
		"success": true,
		"data": {
			"id": "space_123",
			"name": "Test Space",
			"domain": "test-domain"
		}
	}`

	client := &retoolsdk.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	space, err := client.GetSpace("space_123")

	assert.NoError(t, err)
	assert.NotNil(t, space)
	assert.Equal(t, "space_123", space.ID)
	assert.Equal(t, "Test Space", space.Name)
	assert.Equal(t, "test-domain", space.Domain)
}

func TestGetSpace_Failure(t *testing.T) {
	response := `{
		"success": false,
		"message": "Space not found"
	}`

	client := &retoolsdk.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 404,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	space, err := client.GetSpace("non_existing_space")

	assert.Error(t, err)
	assert.Nil(t, space)
	assert.EqualError(t, err, "Space not found")
}

func TestListSpaces_Success(t *testing.T) {
	mockResponse := `{
		"success": true,
		"data": [
			{
				"id": "space_123",
				"name": "Test Space 1",
				"domain": "test-domain1"
			},
			{
				"id": "space_456",
				"name": "Test Space 2",
				"domain": "test-domain2"
			}
		],
		"has_more": false
	}`

	client := &retoolsdk.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(mockResponse))),
				},
			},
		},
	}

	spaces, err := client.ListSpaces()

	assert.NoError(t, err)
	assert.Len(t, spaces, 2)
	assert.Equal(t, "space_123", spaces[0].ID)
	assert.Equal(t, "space_456", spaces[1].ID)
}

func TestListSpaces_Failure(t *testing.T) {
	response := `{
		"success": false,
		"message": "Unauthorized"
	}`

	client := &retoolsdk.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 403,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	spaces, err := client.ListSpaces()

	assert.Error(t, err)
	assert.Nil(t, spaces)
	assert.EqualError(t, err, "Unauthorized")
}

func TestUpdateSpace_Success(t *testing.T) {
	response := `{
		"success": true,
		"data": {
			"id": "space_123",
			"name": "Updated Space",
			"domain": "updated-domain"
		}
	}`

	client := &retoolsdk.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	space, err := client.UpdateSpace("space_123", "Updated Space", "updated-domain")

	assert.NoError(t, err)
	assert.NotNil(t, space)
	assert.Equal(t, "Updated Space", space.Name)
	assert.Equal(t, "updated-domain", space.Domain)
}

func TestUpdateSpace_Failure(t *testing.T) {
	response := `{
		"success": false,
		"message": "Failed to update space"
	}`

	client := &retoolsdk.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 400,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	space, err := client.UpdateSpace("space_123", "Invalid Space", "")

	assert.Error(t, err)
	assert.Nil(t, space)
	assert.EqualError(t, err, "Failed to update space")
}

func TestCreateSpace_Success(t *testing.T) {
	response := `{
		"success": true,
		"data": {
			"id": "space_789",
			"name": "New Space",
			"domain": "new-domain"
		}
	}`

	client := &retoolsdk.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	options := &retoolsdk.CreateSpaceOptions{
		CopySSOSettings:              true,
		CopyBrandingAndThemeSettings: true,
		UsersToCopyAsAdmins:          []string{"admin_1"},
		CreateAdminUser:              true,
	}

	space, err := client.CreateSpace("New Space", "new-domain", options)

	assert.NoError(t, err)
	assert.NotNil(t, space)
	assert.Equal(t, "New Space", space.Name)
	assert.Equal(t, "new-domain", space.Domain)
}

func TestCreateSpace_Failure(t *testing.T) {
	response := `{
		"success": false,
		"message": "Invalid domain name"
	}`

	client := &retoolsdk.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 400,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	options := &retoolsdk.CreateSpaceOptions{
		CopySSOSettings:              true,
		CopyBrandingAndThemeSettings: false,
	}

	space, err := client.CreateSpace("Test Space", "invalid_domain", options)

	assert.Error(t, err)
	assert.Nil(t, space)
	assert.EqualError(t, err, "Invalid domain name")
}

func TestDeleteSpace_Success(t *testing.T) {
	response := `{
		"success": true
	}`

	client := &retoolsdk.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	err := client.DeleteSpace("space_123")

	assert.NoError(t, err)
}

func TestDeleteSpace_Failure(t *testing.T) {
	response := `{
		"success": false,
		"message": "Space not found"
	}`

	client := &retoolsdk.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 404,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	err := client.DeleteSpace("non_existing_space")

	assert.Error(t, err)
	assert.EqualError(t, err, "Space not found")
}
