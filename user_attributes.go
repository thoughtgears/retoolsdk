package retoolsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// UserAttribute is a struct that contains the name and value of a user attribute.
type UserAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type OrganizationAttribute struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Label                 string `json:"label"`
	DataType              string `json:"data_type"`
	DefaultValue          string `json:"default_value"`
	IntercomAttributeName string `json:"intercom_attribute_name"`
}

// UpdateUserAttributes Available from API version 2.1.0+ and onprem version 3.20.1+.
// Adds or updates a user attribute, and returns the updated user metadata. The API token must have the "Users > Write" scope.
func (c *Client) UpdateUserAttributes(id string, attributes []UserAttribute) (map[string]interface{}, error) {
	if len(attributes) == 0 {
		return nil, errors.New("no attributes provided")
	}

	requestBody, err := json.Marshal(attributes)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/users/%s/user_attributes", c.BaseURL, id), bytes.NewBuffer(requestBody))

	var response Response[map[string]interface{}]

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if !response.Success {
		return nil, errors.New(response.Message)
	}

	metadata := make(map[string]interface{})

	if innerMetadata, ok := (response.Data)["metadata"].(map[string]interface{}); ok {
		for key, value := range innerMetadata {
			metadata[key] = value
		}
	}

	return metadata, nil
}

// DeleteUserAttribute Available from API version 2.1.0+ and onprem version 3.20.1+.
// Deletes a user attribute, and returns the updated user metadata. The API token must have the "Users > Write" scope.
func (c *Client) DeleteUserAttribute(id, attribute string) (interface{}, error) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%s/user_attributes/%s", c.BaseURL, id, attribute), nil)

	var response Response[User]

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if !response.Success {
		return nil, errors.New(response.Message)
	}

	if response.Data != (User{}) {
		return response.Data.Metadata, nil
	}

	return nil, nil
}

// GetOrganizationAttributes gets the list of currently configured user attributes for the organization.
// The API token must have the "Users > Read" scope.
func (c *Client) GetOrganizationAttributes() ([]OrganizationAttribute, error) {
	baseURL := fmt.Sprintf("%s/user_attributes", c.BaseURL)

	return doPaginatedRequest[OrganizationAttribute](c, "GET", baseURL, nil, url.Values{})
}
