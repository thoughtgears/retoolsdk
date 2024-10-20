package retoolsdk_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	retool "github.com/thoughtgears/retoolsdk"
	"io"
	"net/http"
	"testing"
)

func TestGetGroup_Success(t *testing.T) {
	response := `
{
	"success": true,
	"data": {
		"id": 123,
		"legacy_id": 123,
		"name": "Test Group",
	    "members": [{
			"id": "user_123",
			"email": "jane.doe@example.com",
			"is_group_admin": false
		}],
		"universal_app_access": true,
		"universal_resource_access": false,
		"universal_workflow_access": false,
		"universal_query_library_access": false,
		"user_list_access": true,
		"audit_log_access": false,
		"unpublished_release_access": false,
		"usage_analytics_access": true,
		"theme_access": false,
		"account_details_access": true,
		"landing_page_app_id": "1eae01b1-49ee-4691-8d4d-b43ef1d7ece4",
		"created_at": "2021-01-01T00:00:00Z",
		"updated_at": "2021-01-01T00:00:00Z",
		"user_invites": [{
			"id": 1,
			"legacy_id": 1,
			"invited_by": "user_321",
			"invited_email": "john.doe@example.com",
			"expires_at": "2021-01-01T00:00:00Z",
			"claimed_by": "user_123",
			"claimed_at": "2021-01-01T00:00:00Z",
			"user_type": "default",
			"metadata": null,
			"created_at": "2021-01-01T00:00:00Z",
			"invite_link": "https://example.com/invite/123"
		}]
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

	group, err := client.GetGroup("1234")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "Test Group", group.Name)
	assert.Equal(t, "jane.doe@example.com", group.Members[0].Email)
	assert.Equal(t, "user_123", group.Members[0].ID)
	assert.Equal(t, "user_321", group.UserInvites[0].InvitedBy)
}

func TestGetGroup_Failure(t *testing.T) {
	response := `
{
	"success": false,
	"message": "Group not found"
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

	group, err := client.GetGroup("1234")
	assert.Error(t, err)
	assert.Nil(t, group)
	assert.Equal(t, "Group not found", err.Error())
}
