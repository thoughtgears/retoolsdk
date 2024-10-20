package retoolsdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// Group is a struct that contains the information about a group.
type Group struct {
	ID                          int          `json:"id,omitempty"`
	LegacyID                    int          `json:"legacy_id,omitempty"`
	Name                        string       `json:"name"`
	Members                     []Member     `json:"members,omitempty"`
	UniversalAppAccess          string       `json:"universal_app_access"`
	UniversalResourceAccess     string       `json:"universal_resource_access"`
	UniversalWorkflowAccess     string       `json:"universal_workflow_access"`
	UniversalQueryLibraryAccess string       `json:"universal_query_library_access"`
	UserInvites                 []UserInvite `json:"user_invites,omitempty"`
	UserListAccess              bool         `json:"user_list_access"`
	AuditLogAccess              bool         `json:"audit_log_access"`
	UnpublishedReleaseAccess    bool         `json:"unpublished_release_access"`
	UsageAnalyticsAccess        bool         `json:"usage_analytics_access"`
	ThemeAccess                 bool         `json:"theme_access"`
	AccountDetailsAccess        bool         `json:"account_details_access"`
	LandingPageAppID            string       `json:"landing_page_app_id"`
	CreatedAt                   string       `json:"created_at,omitempty"`
	UpdatedAt                   string       `json:"updated_at,omitempty"`
}

// Member is a struct that contains the information about a group member.
type Member struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	IsGroupAdmin bool   `json:"is_group_admin"`
}

// UserInvite is a struct that contains the information about a member invite to the group.
type UserInvite struct {
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
	return doSingleRequest[Group](c, "GET", baseURL, nil)
}

// ListGroups get all permission groups for an organization or space. The API token must have the "Groups > Read" scope.
func (c *Client) ListGroups() ([]Group, error) {
	baseURL := fmt.Sprintf("%s/groups", c.BaseURL)
	return doPaginatedRequest[Group](c, baseURL, url.Values{})
}

// Constants for the allowed access levels.
const (
	NoneAccess = "none"
	UseAccess  = "use"
	EditAccess = "edit"
	OwnAccess  = "own"
)

// Validate ensures that the options provided in CreateGroupOpts have valid values.
func (g *Group) Validate() error {
	validAccessLevels := map[string]struct{}{
		NoneAccess: {}, UseAccess: {}, EditAccess: {}, OwnAccess: {},
	}

	validQueryLevels := map[string]struct{}{
		NoneAccess: {}, UseAccess: {}, EditAccess: {},
	}

	if _, ok := validAccessLevels[g.UniversalAppAccess]; !ok && g.UniversalAppAccess != "" {
		return fmt.Errorf("invalid value for UniversalAppAccess: %s", g.UniversalAppAccess)
	}

	if _, ok := validAccessLevels[g.UniversalResourceAccess]; !ok && g.UniversalResourceAccess != "" {
		return fmt.Errorf("invalid value for UniversalResourceAccess: %s", g.UniversalResourceAccess)
	}

	if _, ok := validAccessLevels[g.UniversalWorkflowAccess]; !ok && g.UniversalWorkflowAccess != "" {
		return fmt.Errorf("invalid value for UniversalWorkflowAccess: %s", g.UniversalWorkflowAccess)
	}

	if _, ok := validQueryLevels[g.UniversalQueryLibraryAccess]; !ok && g.UniversalQueryLibraryAccess != "" {
		return fmt.Errorf("invalid value for UniversalQueryLibraryAccess: %s", g.UniversalQueryLibraryAccess)
	}

	return nil
}

// CreateGroup creates a group and returns the created group. The API token must have the "Groups > Write" scope.
func (c *Client) CreateGroup(group *Group) (*Group, error) {
	if err := group.Validate(); err != nil {
		return nil, err
	}

	baseURL := fmt.Sprintf("%s/groups", c.BaseURL)
	return doSingleRequest[Group](c, "POST", baseURL, group)
}

// UpdateGroup update a group in an organization using JSON Patch (RFC 6902). Returns the updated group. The API token must have the "Groups > Write" scope.
func (c *Client) UpdateGroup(id string, operations []UpdateOperations) (*Group, error) {
	if len(operations) == 0 {
		return nil, errors.New("no operations provided")
	}

	for _, op := range operations {
		if err := op.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed for operation: %w", err)
		}
	}

	type body struct {
		Operations []UpdateOperations `json:"operations"`
	}

	var requestBody body
	requestBody.Operations = operations

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	baseURL := fmt.Sprintf("%s/groups/%s", c.BaseURL, id)
	return doSingleRequest[Group](c, "PATCH", baseURL, requestBodyJSON)
}
