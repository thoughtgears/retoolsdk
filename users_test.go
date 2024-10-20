package retoolsdk_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	retool "github.com/thoughtgears/retoolsdk"

	"github.com/stretchr/testify/assert"
)

func TestGetUser_Success(t *testing.T) {
	response := `
{
	"success": true,
	"data": {
		"id": "user_123",
		"legacy_id": 123,
		"email": "jane.doe@example.com",
		"active": true,
		"created_at": "2021-01-01T00:00:00Z",
		"last_active": "2021-01-01T00:00:00Z",
		"first_name": "Jane",
		"last_name": "Doe",
		"metadata": null,
		"is_admin": false,
		"user_type": "user"
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

	user, err := client.GetUser("user_123")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user_123", user.ID)
	assert.Equal(t, "jane.doe@example.com", user.Email)
}

func TestGetUser_EmptyData(t *testing.T) {
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
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	user, err := client.GetUser("user_123")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUser_MisformattedID(t *testing.T) {
	response := `
{
	"success": false,
	"message": "User sid is misformatted: userId"
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

	user, err := client.GetUser("user_123")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestListUser_Success(t *testing.T) {
	response := `
{
	"success": true,
	"data": [{
		"id": "user_123",
		"legacy_id": 123,
		"email": "jane.doe@example.com",
		"active": true,
		"created_at": "2021-01-01T00:00:00Z",
		"last_active": "2021-01-01T00:00:00Z",
		"first_name": "Jane",
		"last_name": "Doe",
		"metadata": null,
		"is_admin": false,
		"user_type": "user"
	}],
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

	users, err := client.ListUsers(nil)
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 1)
	assert.Equal(t, "user_123", users[0].ID)
	assert.Equal(t, "jane.doe@example.com", users[0].Email)
}

func TestListUsers_Pagination(t *testing.T) {
	mockResponsePage1 := `{
		"success": true,
		"data": [{
			"id": "user_123",
			"email": "jane.doe@example.com",
			"first_name": "Jane",
			"last_name": "Doe"
		}],
		"next_token": "next_token",
		"has_more": true
	}`

	mockResponsePage2 := `{
		"success": true,
		"data": [{
			"id": "user_456",
			"email": "john.doe@example.com",
			"first_name": "John",
			"last_name": "Doe"
		}],
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

	users, err := client.ListUsers(nil)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, "user_123", users[0].ID)
	assert.Equal(t, "user_456", users[1].ID)
}

func TestListUsers_EmptyResponse(t *testing.T) {
	response := `{
		"success": true,
		"data": [],
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

	users, err := client.ListUsers(nil)

	assert.NoError(t, err)
	assert.Nil(t, users)
	assert.Len(t, users, 0)
}

func TestCreateUser_Success(t *testing.T) {
	response := `
{
	"success": true,
	"data": {
		"id": "user_123",
		"email": "jane.doe@example.com",
		"first_name": "Jane",
		"last_name": "Doe",
		"user_type": "default",
		"active": true
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

	opts := &retool.CreateUserOpts{
		Active: true,
		Type:   retool.UserTypeDefault,
	}

	user, err := client.CreateUser("jane.doe@example.com", "Jane", "Doe", opts)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user_123", user.ID)
	assert.Equal(t, "jane.doe@example.com", user.Email)
	assert.Equal(t, "Jane", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "default", user.UserType)
	assert.True(t, user.Active)
}

func TestCreateUser_UserAlreadyExists(t *testing.T) {
	response := `
{
	"success": false,
	"message": "User with email jane.doe@example.com already exists"
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

	opts := &retool.CreateUserOpts{
		Active: true,
		Type:   retool.UserTypeDefault,
	}

	user, err := client.CreateUser("jane.doe@example.com", "Jane", "Doe", opts)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "User with email jane.doe@example.com already exists", err.Error())
}

func TestCreateUser_InvalidEmailType(t *testing.T) {
	response := `
{
	"success": false,
	"message": "Invalid email type for body parameter: email"
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

	opts := &retool.CreateUserOpts{
		Active: true,
		Type:   retool.UserTypeDefault,
	}

	user, err := client.CreateUser("invalid-email", "Jane", "Doe", opts)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "Invalid email type for body parameter: email", err.Error())
}

func TestCreateUser_InternalServerError(t *testing.T) {
	response := `
{
	"success": false,
	"message": "Internal server error"
}`

	client := &retool.Client{
		BaseURL: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &MockTransport{
				Response: &http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(response))),
				},
			},
		},
	}

	opts := &retool.CreateUserOpts{
		Active: true,
		Type:   retool.UserTypeDefault,
	}

	user, err := client.CreateUser("jane.doe@example.com", "Jane", "Doe", opts)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "Internal server error", err.Error())
}

func TestUpdateUser_Success(t *testing.T) {
	response := `
{
	"success": true,
	"data": {
		"id": "123",
		"email": "john.doe@example.com",
		"first_name": "NewFirstName",
		"last_name": "Doe"
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

	operations := []retool.UpdateUserOperations{
		{
			Op:    "replace",
			Path:  "first_name",
			Value: "NewFirstName",
		},
	}

	updatedUser, err := client.UpdateUser("123", operations)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "123", updatedUser.ID)
	assert.Equal(t, "john.doe@example.com", updatedUser.Email)
	assert.Equal(t, "NewFirstName", updatedUser.FirstName)
	assert.Equal(t, "Doe", updatedUser.LastName)
}

func TestUpdateUser_Failure(t *testing.T) {
	response := `
{
 "success": false,
 "message": "Patched document failed schema validation: String must contain at least 1 character(s): first_name"
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

	operations := []retool.UpdateUserOperations{
		{
			Op:    "replace",
			Path:  "first_name",
			Value: "",
		},
	}

	updatedUser, err := client.UpdateUser("123", operations)
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Equal(t, "Patched document failed schema validation: String must contain at least 1 character(s): first_name", err.Error())
}

func TestDeleteUser_Success(t *testing.T) {
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

	err := client.DeleteUser("user_123")
	assert.NoError(t, err)
}

func TestDeleteUser_Failure(t *testing.T) {
	response := `
{
	"success": false,
	"message": "No user 'user_123' found for organization"
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

	err := client.DeleteUser("user_123")
	assert.Error(t, err)
	assert.Equal(t, "No user 'user_123' found for organization", err.Error())
}
