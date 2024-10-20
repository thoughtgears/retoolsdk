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
		"universal_app_access": "none",
		"universal_resource_access": "none",
		"universal_workflow_access": "none",
		"universal_query_library_access": "none",
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

func TestCreateGroup_Success(t *testing.T) {
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
		"universal_app_access": "none",
		"universal_resource_access": "none",
		"universal_workflow_access": "none",
		"universal_query_library_access": "none",
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

	requestGroup := &retool.Group{
		Name:                        "Test Group",
		UniversalAppAccess:          retool.NoneAccess,
		UniversalResourceAccess:     retool.NoneAccess,
		UniversalWorkflowAccess:     retool.NoneAccess,
		UniversalQueryLibraryAccess: retool.NoneAccess,
		UserListAccess:              true,
		AuditLogAccess:              false,
		UnpublishedReleaseAccess:    false,
		UsageAnalyticsAccess:        true,
		ThemeAccess:                 false,
		AccountDetailsAccess:        true,
		LandingPageAppID:            "1eae01b1-49ee-4691-8d4d-b43ef1d7ece4",
		Members: []retool.Member{
			{
				ID:           "user_123",
				IsGroupAdmin: false,
			},
		},
		UserInvites: []retool.UserInvite{
			{
				ID:           1,
				LegacyID:     1,
				InvitedBy:    "user_321",
				InvitedEmail: "john.doe@example.com",
				ExpiresAt:    "2021-01-01T00:00:00Z",
				ClaimedBy:    "user_123",
				ClaimedAt:    "2021-01-01T00:00:00Z",
				UserType:     "default",
				InviteLink:   "https://example.com/invite/123",
			},
		},
	}

	group, err := client.CreateGroup(requestGroup)

	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "Test Group", group.Name)
	assert.Equal(t, "jane.doe@example.com", group.Members[0].Email)
	assert.Equal(t, "user_123", group.Members[0].ID)
	assert.Equal(t, "user_321", group.UserInvites[0].InvitedBy)
}

func TestCreateGroup_ValidationFailure(t *testing.T) {
	requestGroup := &retool.Group{
		UniversalAppAccess:      "invalid_access",
		UniversalResourceAccess: "another_invalid",
	}

	client := &retool.Client{
		BaseURL:    "https://example.com",
		HTTPClient: &http.Client{},
	}

	_, err := client.CreateGroup(requestGroup)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid value for UniversalAppAccess")
}

func TestCreateGroup_Failure(t *testing.T) {
	response := `
{
	"success": false,
	"message": "some error message"
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

	requestGroup := &retool.Group{}

	group, err := client.CreateGroup(requestGroup)

	assert.Error(t, err)
	assert.Nil(t, group)
}
