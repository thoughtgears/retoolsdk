package retoolsdk_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	retool "github.com/thoughtgears/retoolsdk"
	"io"
	"net/http"
	"testing"
)

func TestGetConfigurationVariable_Success(t *testing.T) {
	response := `{
		"success": true,
		"data": {
			"id": "config_var_123",
			"name": "Test Variable",
			"description": "A test variable",
			"secret": false,
			"values": [
				{
					"key": "development",
					"value": "dev_value"
				}
			]
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

	configVar, err := client.GetConfigurationVariable("config_var_123")
	assert.NoError(t, err)
	assert.NotNil(t, configVar)
	assert.Equal(t, "Test Variable", configVar.Name)
	assert.Equal(t, "dev_value", configVar.Values[0].Value)
}

func TestGetConfigurationVariable_Failure(t *testing.T) {
	response := `{
		"success": false,
		"message": "Configuration variable not found"
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

	configVar, err := client.GetConfigurationVariable("config_var_123")
	assert.Error(t, err)
	assert.Nil(t, configVar)
}

func TestListConfigurationVariables_Success(t *testing.T) {
	response := `{
		"success": true,
		"data": [
			{
				"id": "config_var_123",
				"name": "Test Variable",
				"description": "A test variable",
				"secret": false,
				"values": [
					{
						"key": "development",
						"value": "dev_value"
					}
				]
			}
		]
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

	configVars, err := client.ListConfigurationVariables()
	assert.NoError(t, err)
	assert.NotNil(t, configVars)
	assert.Len(t, configVars, 1)
	assert.Equal(t, "Test Variable", configVars[0].Name)
	assert.Equal(t, "dev_value", configVars[0].Values[0].Value)
}

func TestListConfigurationVariables_Failure(t *testing.T) {
	response := `{
		"success": false,
		"message": "Configuration variables not found"
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

	configVars, err := client.ListConfigurationVariables()
	assert.Error(t, err)
	assert.Nil(t, configVars)
}

func TestCreateConfigurationVariable_Success(t *testing.T) {
	response := `{
		"success": true,
		"data": {
			"id": "config_var_123",
			"name": "Test Variable",
			"description": "A test variable",
			"secret": false,
			"values": [
				{
					"key": "development",
					"value": "prod_value"
				}
			]
		}
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

	configVar, err := client.CreateConfigurationVariable("Test Variable", "A test variable", false, []retool.Value{{EnvironmentId: "production", Value: "prod_value"}})

	assert.NoError(t, err)
	assert.NotNil(t, configVar)
	assert.Equal(t, "Test Variable", configVar.Name)
	assert.Equal(t, "prod_value", configVar.Values[0].Value)
}

func TestCreateConfigurationVariable_Failure(t *testing.T) {
	response := `{
		"success": false,
		"message": "Failed to create configuration variable"
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

	configVar, err := client.CreateConfigurationVariable("Test Variable", "A test variable", false, []retool.Value{{EnvironmentId: "production", Value: "prod_value"}})
	assert.Error(t, err)
	assert.Nil(t, configVar)
}

func TestUpdateConfigurationVariable_Success(t *testing.T) {
	response := `{
		"success": true,
		"data": {
			"id": "config_var_123",
			"name": "Test Variable",
			"description": "A test variable",
			"secret": false,
			"values": [
				{
					"key": "development",
					"value": "prod_value"
				}
			]
		}
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

	configVar, err := client.UpdateConfigurationVariable("config_var_123", "Test Variable", "A test variable", false, []retool.Value{{EnvironmentId: "production", Value: "prod_value"}})
	assert.NoError(t, err)
	assert.NotNil(t, configVar)
	assert.Equal(t, "Test Variable", configVar.Name)
}

func TestUpdateConfigurationVariable_Failure(t *testing.T) {
	response := `{
		"success": false,
		"message": "Failed to update configuration variable"
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

	configVar, err := client.UpdateConfigurationVariable("config_var_123", "Test Variable", "A test variable", false, []retool.Value{{EnvironmentId: "production", Value: "prod_value"}})
	assert.Error(t, err)
	assert.Nil(t, configVar)
}

func TestClient_DeleteConfigurationVariable_Success(t *testing.T) {
	client := &retool.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 204,
					Body:       nil,
				},
			},
		},
	}

	err := client.DeleteConfigurationVariable("config_var_123")
	assert.NoError(t, err)
}

func TestClient_DeleteConfigurationVariable_Failure(t *testing.T) {
	response := `{
		"success": false,
		"message": "Failed to delete configuration variable"
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

	err := client.DeleteConfigurationVariable("config_var_123")
	assert.Error(t, err)
}
