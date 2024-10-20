package retool_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type User struct {
	ID         string      `json:"id"`
	LegacyID   int         `json:"legacy_id"`
	Email      string      `json:"email"`
	Active     bool        `json:"active"`
	CreatedAt  string      `json:"created_at"`
	LastActive string      `json:"last_active"`
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Metadata   interface{} `json:"metadata"`
	IsAdmin    bool        `json:"is_admin"`
	UserType   string      `json:"user_type"`
}

// GetUser returns a user by ID, or nil if the user is not found.
// If the user ID is misformatted, an error is returned.
func (c *Client) GetUser(id string) (*User, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%s", c.BaseURL, id), nil)

	var response Response[User]

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if !response.Success && response.Message == "User sid is misformatted: userId" {
		return nil, errors.New(response.Message)
	}

	if !response.Success == false && response.Message == "User not found" {
		return nil, nil
	}

	return response.Data.Single, nil
}

// ListUserOpts is a struct that contains optional query parameters for ListUsers
type ListUserOpts struct {
	Email     string
	FirstName string
	LastName  string
}

// ListUsers returns a list of users. If there are no users, an empty list is returned.
// ListUsers can take email, first or last name as a query parameter to filter the list.
func (c *Client) ListUsers(opts *ListUserOpts) ([]User, error) {
	baseURL := fmt.Sprintf("%s/users", c.BaseURL)

	query := make(url.Values)

	if opts != nil {
		if opts.Email != "" {
			query.Add("email", opts.Email)
		}
		if opts.FirstName != "" {
			query.Add("first_name", opts.FirstName)
		}
		if opts.LastName != "" {
			query.Add("last_name", opts.LastName)
		}

		if len(query) > 0 {
			baseURL = fmt.Sprintf("%s?%s", baseURL, query.Encode())
		}
	}

	var allUsers []User
	var nextToken string
	hasMore := true

	for hasMore {
		if nextToken != "" {
			query.Set("next", nextToken)
		}

		urlWithQuery := fmt.Sprintf("%s?%s", baseURL, query.Encode())

		req, err := http.NewRequest("GET", urlWithQuery, nil)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		var response Response[User]

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("making request: %w", err)
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("decoding response: %w", err)
		}

		allUsers = append(allUsers, response.Data.List...)

		nextToken = response.NextToken
		hasMore = response.HasMore
	}

	return allUsers, nil
}
