package retoolsdk_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	retool "github.com/thoughtgears/retoolsdk"
	"io"
	"net/http"
	"testing"
)

func TestUpdateUserAttributes_Success(t *testing.T) {
	response := `
	{
		"success": true,
		"data": {
			"metadata": {
				"attribute1": "value1",
				"attribute2": "value2"
			}
		}
	}`

	client := &retool.Client{
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

	attributes := []retool.UserAttribute{
		{Name: "attribute1", Value: "value1"},
		{Name: "attribute2", Value: "value2"},
	}

	expectedMetadata := map[string]interface{}{
		"attribute1": "value1",
		"attribute2": "value2",
	}

	metadata, err := client.UpdateUserAttributes("user_123", attributes)
	assert.NoError(t, err)
	assert.NotNil(t, metadata)
	assert.Equal(t, expectedMetadata, metadata)
}

func TestUserAttributes_Failure(t *testing.T) {
	response := `
	{
		"success": false,
		"message": "Invalid attributes"
	}`

	client := &retool.Client{
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

	attributes := []retool.UserAttribute{
		{Name: "attribute1", Value: "value1"},
	}

	metadata, err := client.UpdateUserAttributes("user_123", attributes)
	assert.Error(t, err)
	assert.Nil(t, metadata)
	assert.Equal(t, "Invalid attributes", err.Error())
}

func TestDeleteUserAttribute_Success(t *testing.T) {
	response := `
	{
		"success": true,
		"data": {
			"id": "user_123",
			"metadata": {
				"attribute1": "value1"
			}
		}
	}`

	client := &retool.Client{
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

	expectedMetadata := map[string]interface{}{
		"attribute1": "value1",
	}

	metadata, err := client.DeleteUserAttribute("user_123", "attribute2")
	assert.NoError(t, err)
	assert.NotNil(t, metadata)
	assert.Equal(t, expectedMetadata, metadata)
}

func TestDeleteUserAttribute_Failure(t *testing.T) {
	response := `
	{
		"success": false,
		"message": "Attribute not found"
	}`

	client := &retool.Client{
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

	metadata, err := client.DeleteUserAttribute("user_123", "attribute1")
	assert.Error(t, err)
	assert.Nil(t, metadata)
	assert.Equal(t, "Attribute not found", err.Error())
}

func TestGetOrganizationAttributes_Success(t *testing.T) {
	response := `
	{
		"success": true,
		"data": [{
					"id": "org_attr_1",
					"name": "attribute1",
					"label": "Attribute 1",
					"data_type": "string",
					"default_value": "default1",
					"intercom_attribute_name": "intercom_attr_1"
				},
				{
					"id": "org_attr_2",
					"name": "attribute2",
					"label": "Attribute 2",
					"data_type": "string",
					"default_value": "default2",
					"intercom_attribute_name": "intercom_attr_2"
				}],
		"total_count": 2,
		"has_more": false
	}`

	client := &retool.Client{
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

	attributes, err := client.GetOrganizationAttributes()
	assert.NoError(t, err)
	assert.NotNil(t, attributes)
	assert.Equal(t, "attribute1", attributes[0].Name)
}

func TestGetOrganizationAttributes_Failure(t *testing.T) {
	response := `
	{
		"success": false,
		"message": "Failed to retrieve attributes"
	}`

	client := &retool.Client{
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

	attributes, err := client.GetOrganizationAttributes()
	assert.Error(t, err)
	assert.Nil(t, attributes)
	assert.Equal(t, "Failed to retrieve attributes", err.Error())
}
