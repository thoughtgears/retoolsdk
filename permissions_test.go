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

func TestValidateAccessLevel_Success(t *testing.T) {
	perm := retool.AccessLevel("own")
	err := perm.Validate()

	assert.NoError(t, err, "Valid access level should not return an error")
}

func TestValidateAccessLevel_Failure(t *testing.T) {
	perm := retool.AccessLevel("invalid")
	err := perm.Validate()

	assert.Error(t, err, "Invalid access level should return an error")
}

func TestValidateObjectType_Success(t *testing.T) {
	perm := retool.ObjectType("app")
	err := perm.Validate()

	assert.NoError(t, err, "Valid object type should not return an error")
}

func TestValidateObjectType_Failure(t *testing.T) {
	perm := retool.ObjectType("invalid")
	err := perm.Validate()

	assert.Error(t, err, "Invalid object type should return an error")
}

func TestGetFolderOrAppAccessList_Success(t *testing.T) {
	response := `
{
  "success": true,
  "data": {
    "group": [
      {
        "subject": {
          "id": "group_123",
          "type": "group"
        },
        "sources": {
          "direct": true,
          "universal": true,
          "groups": [
            {
              "id": "group_456",
              "type": "group"
            }
          ],
          "inherited": {
            "id": "app_789",
            "type": "app"
          }
        },
        "accessLevel": "own"
      }
    ],
    "user": [
      {
        "subject": {
          "id": "user_123",
          "type": "user"
        },
        "sources": {
          "direct": true,
          "universal": true,
          "groups": [
            {
              "id": "group_456",
              "type": "group"
            }
          ],
          "inherited": {
            "id": "app_789",
            "type": "app"
          }
        },
        "accessLevel": "edit"
      }
    ],
    "userInvite": [
      {
        "subject": {
          "id": "invite_123",
          "type": "userInvite"
        },
        "sources": {
          "direct": true,
          "universal": true,
          "groups": [
            {
              "id": "group_789",
              "type": "group"
            }
          ],
          "inherited": {
            "id": "app_123",
            "type": "app"
          }
        },
        "accessLevel": "use"
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

	accessData, err := client.GetFolderOrAppAccessList("123", "app")

	assert.NoError(t, err)
	assert.NotNil(t, accessData)
	assert.Len(t, accessData.Group, 1)
	assert.Len(t, accessData.User, 1)
	assert.Len(t, accessData.UserInvite, 1)
	assert.Equal(t, "group_123", accessData.Group[0].Subject.ID)
}

func TestGetFolderOrAppAccessList_Failure(t *testing.T) {
	response := `
{
	"success": false,
	"message": "Object not found"
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

	accessData, err := client.GetFolderOrAppAccessList("123", "app")

	assert.Error(t, err)
	assert.Nil(t, accessData)
	assert.Equal(t, "Object not found", err.Error())
}

func TestListGroupObjectPermissions_Success(t *testing.T) {
	response := `
{
  "success": true,
  "data": [
    {
      "type": "folder",
      "id": "123",
      "access_level": "own"
    }
  ],
  "total_count": 1,
  "next_token": "",
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

	resp, err := client.ListGroupObjectPermissions("group", "folder", 123)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "folder", resp[0].Type)
	assert.Equal(t, "123", resp[0].ID)
	assert.Equal(t, "own", resp[0].AccessLevel)
}

func TestListGroupObjectPermissions_Pagination(t *testing.T) {
	mockResponsePage1 := `
{
  "success": true,
  "data": [
    {
      "type": "folder",
      "id": "123",
      "access_level": "own"
    }
  ],
  "total_count": 2,
  "next_token": "next_token",
  "has_more": true
}`

	mockResponsePage2 := `
{
  "success": true,
  "data": [
    {
      "type": "folder",
      "id": "321",
      "access_level": "own"
    }
  ],
  "total_count": 2,
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

	resp, err := client.ListGroupObjectPermissions("group", "folder", 123)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp, 2)
	assert.Equal(t, "123", resp[0].ID)
	assert.Equal(t, "321", resp[1].ID)
}

func TestListGroupObjectPermissions_Failure(t *testing.T) {
	response := `
{
	"success": false,
	"message": "Object not found"
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

	resp, err := client.ListGroupObjectPermissions("group", "folder", 123)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "Object not found", err.Error())
}

func TestGrantPermission_Success(t *testing.T) {
	response := `
	{
		"success": true,
		"data": [
			{
				"id": "user_123",
				"type": "user",
				"access_level": "own"
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

	subjects, err := client.GrantPermission("user", "user_123", retool.ObjectType("app"), "app_123", retool.AccessLevel("own"))
	assert.NoError(t, err)
	assert.NotNil(t, subjects)
	assert.Equal(t, "user_123", subjects[0].ID)
	assert.Equal(t, "user", subjects[0].Type)
	assert.Equal(t, "own", subjects[0].AccessLevel)
}

func TestGrantPermission_Failure(t *testing.T) {
	response := `
	{
		"success": false,
		"message": "Invalid access level"
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

	subjects, err := client.GrantPermission("user", "user_123", retool.ObjectType("app"), "app_123", retool.AccessLevel("invalid"))
	assert.Error(t, err)
	assert.Nil(t, subjects)
	assert.Equal(t, "Invalid access level", err.Error())
}

func TestGrantPermission_Pagination(t *testing.T) {
	mockResponsePage1 := `
	{
		"success": true,
		"data": [
			{
				"id": "user_123",
				"type": "user",
				"access_level": "own"
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
				"id": "user_456",
				"type": "user",
				"access_level": "edit"
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

	subjects, err := client.GrantPermission("user", "user_123", retool.ObjectType("app"), "app_123", retool.AccessLevel("own"))
	assert.NoError(t, err)
	assert.NotNil(t, subjects)
	assert.Len(t, subjects, 2)
	assert.Equal(t, "user_123", subjects[0].ID)
	assert.Equal(t, "user_456", subjects[1].ID)
}

func TestRevokePermission_Success(t *testing.T) {
	response := `
	{
		"success": true,
		"data": [
			{
				"id": "user_123",
				"type": "user",
				"access_level": "own"
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

	subjects, err := client.RevokePermission("user", "user_123", retool.ObjectType("app"), "app_123")
	assert.NoError(t, err)
	assert.NotNil(t, subjects)
	assert.Equal(t, "user_123", subjects[0].ID)
	assert.Equal(t, "user", subjects[0].Type)
	assert.Equal(t, "own", subjects[0].AccessLevel)
}

func TestRevokePermission_Failure(t *testing.T) {
	response := `
	{
		"success": false,
		"message": "User not found"
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

	subjects, err := client.RevokePermission("user", "user_123", retool.ObjectType("app"), "app_123")
	assert.Error(t, err)
	assert.Nil(t, subjects)
	assert.Equal(t, "User not found", err.Error())
}

func TestRevokePermission_Pagination(t *testing.T) {
	mockResponsePage1 := `
	{
		"success": true,
		"data": [
			{
				"id": "user_123",
				"type": "user",
				"access_level": "own"
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
				"id": "user_456",
				"type": "user",
				"access_level": "edit"
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

	subjects, err := client.RevokePermission("user", "user_123", retool.ObjectType("app"), "app_123")
	assert.NoError(t, err)
	assert.NotNil(t, subjects)
	assert.Len(t, subjects, 2)
	assert.Equal(t, "user_123", subjects[0].ID)
	assert.Equal(t, "user_456", subjects[1].ID)
}
