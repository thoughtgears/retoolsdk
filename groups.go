package retoolsdk

import (
	"fmt"
	"net/url"
)

// Group is a struct that contains the information about a group.
type Group struct {
	ID                          int      `json:"id,omitempty"`
	LegacyID                    int      `json:"legacy_id,omitempty"`
	Name                        string   `json:"name"`
	Members                     []Member `json:"members,omitempty"`
	UniversalAppAccess          bool     `json:"universal_app_access"`
	UniversalResourceAccess     bool     `json:"universal_resource_access"`
	UniversalWorkflowAccess     bool     `json:"universal_workflow_access"`
	UniversalQueryLibraryAccess bool     `json:"universal_query_library_access"`
	UserInvites                 []Invite `json:"user_invites,omitempty"`
	UserListAccess              bool     `json:"user_list_access"`
	AuditLogAccess              bool     `json:"audit_log_access"`
	UnpublishedReleaseAccess    bool     `json:"unpublished_release_access"`
	UsageAnalyticsAccess        bool     `json:"usage_analytics_access"`
	ThemeAccess                 bool     `json:"theme_access"`
	AccountDetailsAccess        bool     `json:"account_details_access"`
	LandingPageAppID            string   `json:"landing_page_app_id"`
	CreatedAt                   string   `json:"created_at"`
	UpdatedAt                   string   `json:"updated_at"`
}

// Member is a struct that contains the information about a group member.
type Member struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	IsGroupAdmin bool   `json:"is_group_admin"`
}

// Invite is a struct that contains the information about a member invite to the group.
type Invite struct {
	ID           int         `json:"id"`
	LegacyID     int         `json:"legacy_id"`
	InvitedBy    string      `json:"invited_by"`
	InvitedEmail string      `json:"invited_email"`
	ExpiresAt    string      `json:"expires_at"`
	ClaimedBy    string      `json:"claimed_by"`
	ClaimedAt    string      `json:"claimed_at"`
	UserType     string      `json:"user_type"`
	Metadata     interface{} `json:"metadata"`
	CreatedAt    string      `json:"created_at"`
	InviteLink   string      `json:"invite_link"`
}

// GetGroup get a group with a given groupId. The API token must have the "Groups > Read" scope.
func (c *Client) GetGroup(id string) (*Group, error) {
	baseURL := fmt.Sprintf("%s/groups/%s", c.BaseURL, id)
	resp, err := c.Do("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	data, err := decodeResponse[Group](resp)
	if err != nil {
		return nil, err
	}

	return &data.Data, nil
}

// ListGroups get all permission groups for an organization or space. The API token must have the "Groups > Read" scope.
func (c *Client) ListGroups() ([]Group, error) {
	baseURL := fmt.Sprintf("%s/groups", c.BaseURL)
	return doPaginatedRequest[Group](c, baseURL, url.Values{})
}
