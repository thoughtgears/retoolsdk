package retoolsdk

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type ConfigurationVariable struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Secret      bool    `json:"secret"`
	Values      []Value `json:"values"`
}

type Value struct {
	EnvironmentId string `json:"environment_id"`
	Value         string `json:"value"`
}

// GetConfigurationVariable available for orgs with configuration variables enabled on Retool Version 3.42+.
// The API token must have the "Configuration Variables > Read" scope.
func (c *Client) GetConfigurationVariable(id string) (*ConfigurationVariable, error) {
	baseURL := fmt.Sprintf("%s/configuration_variables/%s", c.BaseURL, id)
	return doSingleRequest[ConfigurationVariable](c, "GET", baseURL, nil)
}

// ListConfigurationVariables available for orgs with configuration variables enabled on Retool Version 3.42+.
// The API token must have the "Configuration Variables > Read" scope.
func (c *Client) ListConfigurationVariables() ([]ConfigurationVariable, error) {
	baseURL := fmt.Sprintf("%s/configuration_variables", c.BaseURL)
	return doPaginatedRequest[ConfigurationVariable](c, "GET", baseURL, nil, url.Values{})
}

// CreateConfigurationVariable available for orgs with configuration variables enabled on Retool Version 3.42+.
// The API token must have the "Configuration Variables > Write" scope.
func (c *Client) CreateConfigurationVariable(name, description string, secret bool, values []Value) (*ConfigurationVariable, error) {
	requestBody := struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Secret      bool    `json:"secret"`
		Values      []Value `json:"values"`
	}{
		Name:        name,
		Description: description,
		Secret:      secret,
		Values:      values,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	baseURL := fmt.Sprintf("%s/configuration_variables", c.BaseURL)
	return doSingleRequest[ConfigurationVariable](c, "POST", baseURL, requestBodyJSON)
}

// UpdateConfigurationVariable update a configuration variable and its values. Available for orgs with configuration
// variables enabled on Retool Version 3.42+. The API token must have the "Configuration Variables > Write" scope.
func (c *Client) UpdateConfigurationVariable(id, name, description string, secret bool, values []Value) (*ConfigurationVariable, error) {
	requestBody := struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Secret      bool    `json:"secret"`
		Values      []Value `json:"values"`
	}{
		Name:        name,
		Description: description,
		Secret:      secret,
		Values:      values,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	baseURL := fmt.Sprintf("%s/configuration_variables/%s", c.BaseURL, id)
	return doSingleRequest[ConfigurationVariable](c, "PUT", baseURL, requestBodyJSON)
}

// DeleteConfigurationVariable deletes a configuration variable and its values. Available for orgs with configuration
// variables enabled on Retool Version 3.42+. The API token must have the "Configuration Variables > Write" scope.
func (c *Client) DeleteConfigurationVariable(id string) error {
	baseURL := fmt.Sprintf("%s/configuration_variables/%s", c.BaseURL, id)
	_, err := doSingleRequest[ConfigurationVariable](c, "DELETE", baseURL, nil)
	return err
}
