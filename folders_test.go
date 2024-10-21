package retoolsdk_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	retool "github.com/thoughtgears/retoolsdk"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFolder_Success(t *testing.T) {
	response := `
	{	
		"success": true,
		"data": {
			"id": "folder_123",
			"legacy_id": "legacy_123",
			"name": "Test Folder",
			"parent_folder_id": "parent_123",
			"is_system_folder": false,
			"folder_type": "app",
			"created_at": "2023-01-01T00:00:00Z",
			"updated_at": "2023-01-01T00:00:00Z"
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

	folder, err := client.GetFolder("folder_123")

	assert.NoError(t, err)
	assert.NotNil(t, folder)
	assert.Equal(t, "folder_123", folder.ID)
}

func TestGetFolder_Failure(t *testing.T) {
	response := `
	{
		"success": false,
		"message": "Folder not found"
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

	folder, err := client.GetFolder("folder_123")

	assert.Error(t, err)
	assert.Nil(t, folder)
	assert.Equal(t, "Folder not found", err.Error())
}

func TestListFolders_Success(t *testing.T) {
	response := `
	{
		"success": true,
		"data": [
			{
				"id": "folder_123",
				"legacy_id": "legacy_123",
				"name": "Test Folder",
				"parent_folder_id": "parent_123",
				"is_system_folder": false,
				"folder_type": "app",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z"
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

	folders, err := client.ListFolders()

	assert.NoError(t, err)
	assert.NotNil(t, folders)
	assert.Len(t, folders, 1)
	assert.Equal(t, "folder_123", folders[0].ID)
}

func TestListFolders_Failure(t *testing.T) {
	response := `
	{
		"success": false,
		"message": "Folders not found"
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

	folders, err := client.ListFolders()

	assert.Error(t, err)
	assert.Nil(t, folders)
	assert.Equal(t, "Folders not found", err.Error())
}

func TestListFolders_Pagination(t *testing.T) {
	mockResponsePage1 := `
	{
		"success": true,
		"data": [
			{
				"id": "folder_123",
				"legacy_id": "legacy_123",
				"name": "Test Folder",
				"parent_folder_id": "parent_123",
				"is_system_folder": false,
				"folder_type": "app",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z"
			}
		],
		"next_token": "next_token",
		"has_more": true
	}`

	mockResponsePage2 := `
	{
		"success": true,
		"data": [
			{
				"id": "folder_456",
				"legacy_id": "legacy_456",
				"name": "Another Folder",
				"parent_folder_id": "parent_456",
				"is_system_folder": false,
				"folder_type": "workflow",
				"created_at": "2023-01-02T00:00:00Z",
				"updated_at": "2023-01-02T00:00:00Z"
			}
		],
		"next_token": "",
		"has_more": false
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("next") == "next_token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, mockResponsePage2)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, mockResponsePage1)
		}
	}))
	defer server.Close()

	client := &retool.Client{
		BaseURL: server.URL,
		HTTPClient: &http.Client{
			Transport: http.DefaultTransport,
		},
	}

	folders, err := client.ListFolders()

	assert.NoError(t, err)
	assert.NotNil(t, folders)
	assert.Len(t, folders, 2)
	assert.Equal(t, "folder_123", folders[0].ID)
	assert.Equal(t, "folder_456", folders[1].ID)
}

func TestCreateFolder_Success(t *testing.T) {
	response := `
	{	
		"success": true,
		"data": {
			"id": "folder_123",
			"legacy_id": "legacy_123",
			"name": "Test Folder",
			"parent_folder_id": "parent_123",
			"is_system_folder": false,
			"folder_type": "app",
			"created_at": "2023-01-01T00:00:00Z",
			"updated_at": "2023-01-01T00:00:00Z"
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

	folder, err := client.CreateFolder("Test Folder", "parent_123", retool.FolderType("app"))

	assert.NoError(t, err)
	assert.NotNil(t, folder)
	assert.Equal(t, "folder_123", folder.ID)
}

func TestCreateFolder_Failure(t *testing.T) {
	response := `
	{
		"success": false,
		"message": "Invalid folder type"
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

	folder, err := client.CreateFolder("Test Folder", "parent_123", retool.FolderType("invalid"))

	assert.Error(t, err)
	assert.Nil(t, folder)
	assert.Equal(t, "validating folder type: invalid value for UserType: invalid", err.Error())
}

func TestUpdateFolder_Success(t *testing.T) {
	response := `
	{	
		"success": true,
		"data": {
			"id": "folder_123",
			"legacy_id": "legacy_123",
			"name": "Updated Folder",
			"parent_folder_id": "parent_123",
			"is_system_folder": false,
			"folder_type": "app",
			"created_at": "2023-01-01T00:00:00Z",
			"updated_at": "2023-01-01T00:00:00Z"
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

	operations := []retool.UpdateOperations{
		{Op: "replace", Path: "/name", Value: "Updated Folder"},
	}

	folder, err := client.UpdateFolder("folder_123", operations)

	assert.NoError(t, err)
	assert.NotNil(t, folder)
	assert.Equal(t, "Updated Folder", folder.Name)
}

func TestUpdateFolder_Failure(t *testing.T) {
	response := `
	{
		"success": false,
		"message": "Folder not found"
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

	operations := []retool.UpdateOperations{
		{Op: "replace", Path: "/name", Value: "Updated Folder"},
	}

	folder, err := client.UpdateFolder("folder_123", operations)

	assert.Error(t, err)
	assert.Nil(t, folder)
	assert.Equal(t, "Folder not found", err.Error())
}

func TestDeleteFolder_Success(t *testing.T) {
	response := `
	{
		"success": true
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

	err := client.DeleteFolder("folder_123")

	assert.NoError(t, err)
}

func TestDeleteFolder_Failure(t *testing.T) {
	response := `
	{
		"success": false,
		"message": "Folder not found"
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

	err := client.DeleteFolder("folder_123")

	assert.Error(t, err)
	assert.Equal(t, "Folder not found", err.Error())
}
