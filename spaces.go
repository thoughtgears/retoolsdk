package retoolsdk

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Space struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetSpace Available for orgs with Spaces enabled. Get space by ID. The API token must have the "Spaces > Read" scope.
func (c *Client) GetSpace(id string) (*Space, error) {
	baseURL := fmt.Sprintf("%s/spaces/%s", c.BaseURL, id)
	return doSingleRequest[Space](c, "GET", baseURL, nil)
}

// ListSpaces Available for orgs with Spaces enabled. List all child spaces of the current space. The API token must have the "Spaces > Read" scope.
func (c *Client) ListSpaces() ([]Space, error) {
	baseURL := fmt.Sprintf("%s/spaces", c.BaseURL)
	return doPaginatedRequest[Space](c, "GET", baseURL, nil, url.Values{})
}

// UpdateSpace Available for orgs with Spaces enabled. Update space by ID. The API token must have the "Spaces > Write" scope.
func (c *Client) UpdateSpace(id, name, domain string) (*Space, error) {
	requestBody := struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
	}{
		Name:   name,
		Domain: domain,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	baseURL := fmt.Sprintf("%s/spaces/%s", c.BaseURL, id)
	return doSingleRequest[Space](c, "PUT", baseURL, requestBodyJSON)
}

// CreateSpaceOptions is a struct that contains optional parameters for CreateSpace.
type CreateSpaceOptions struct {
	CopySSOSettings              bool     `json:"copy_sso_settings"`
	CopyBrandingAndThemeSettings bool     `json:"copy_branding_and_theme_settings"`
	UsersToCopyAsAdmins          []string `json:"users_to_copy_as_admins"`
	CreateAdminUser              bool     `json:"create_admin_user"`
}

// CreateSpace Available for orgs with Spaces enabled. Creates a new child space and returns it. The API token must have the "Spaces > Write" scope.
func (c *Client) CreateSpace(name, domain string, options *CreateSpaceOptions) (*Space, error) {
	requestBody := struct {
		Name    string              `json:"name"`
		Domain  string              `json:"domain"`
		Options *CreateSpaceOptions `json:"options,omitempty"`
	}{
		Name:    name,
		Domain:  domain,
		Options: options,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	baseURL := fmt.Sprintf("%s/spaces", c.BaseURL)
	return doSingleRequest[Space](c, "POST", baseURL, requestBodyJSON)
}

// DeleteSpace Available for orgs with Spaces enabled. Delete a space by ID. The API token must have the "Spaces > Write" scope.
func (c *Client) DeleteSpace(id string) error {
	baseURL := fmt.Sprintf("%s/spaces/%s", c.BaseURL, id)
	_, err := doSingleRequest[any](c, "DELETE", baseURL, nil)
	return err
}
