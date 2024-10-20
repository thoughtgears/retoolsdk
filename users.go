package retoolsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type User struct {
	ID         string      `json:"id,omitempty"`
	LegacyID   int         `json:"legacy_id,omitempty"`
	Email      string      `json:"email"`
	Active     bool        `json:"active,omitempty"`
	CreatedAt  string      `json:"created_at,omitempty"`
	LastActive string      `json:"last_active,omitempty"`
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Metadata   interface{} `json:"metadata,omitempty"`
	IsAdmin    bool        `json:"is_admin,omitempty"`
	UserType   string      `json:"user_type,omitempty"`
}

// GetUser returns the user. The API token must have the "Users > Read" scope.
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

	if response.Success {
		return response.Data.Single, nil
	}

	if !response.Success && response.Message == "User sid is misformatted: userId" {
		return nil, errors.New(response.Message)
	}

	return nil, nil
}

// ListUserOpts is a struct that contains optional query parameters for ListUsers.
type ListUserOpts struct {
	Email     string
	FirstName string
	LastName  string
}

// ListUsers returns a list of users. The API token must have the "Users > Read" scope.
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

// CreateUserOpts is a struct that contains optional parameters for CreateUser.
type CreateUserOpts struct {
	Active bool
	Type   string
}

// User types allowed.
const (
	UserTypeDefault = "default"
	UserTypeEmbed   = "embed"
	UserTypeMobile  = "mobile"
)

// CreateUser creates a user and returns the created user. The API token must have the "Users > Write" scope.
func (c *Client) CreateUser(email, firstName, lastName string, opts *CreateUserOpts) (*User, error) {
	newUser := &User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	if opts != nil {
		if opts.Type == "" {
			newUser.UserType = UserTypeDefault
		} else if opts.Type != UserTypeDefault && opts.Type != UserTypeEmbed && opts.Type != UserTypeMobile {
			return nil, errors.New("invalid newUser type: must be 'default', 'embed', or 'mobile'")
		} else {
			newUser.UserType = opts.Type
		}

		newUser.Active = opts.Active
	} else {
		newUser.UserType = UserTypeDefault
		newUser.Active = false
	}

	newUserJSON, err := json.Marshal(newUser)
	if err != nil {
		return nil, fmt.Errorf("marshalling newUser: %w", err)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/users", c.BaseURL), bytes.NewBuffer(newUserJSON))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	var response Response[User]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if !response.Success && response.Message == fmt.Sprintf("User with email %s already exists", email) {
		return nil, errors.New(response.Message)
	}

	if !response.Success && response.Message == "Invalid email type for body parameter: email" {
		return nil, errors.New(response.Message)
	}

	if !response.Success && response.Message == "Internal server error" {
		return nil, errors.New(response.Message)
	}

	return response.Data.Single, nil
}

// UpdateUserOperations is a struct that contains the operations to update a user.
type UpdateUserOperations struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// User operations allowed.
const (
	UserOpAdd     = "add"
	UserOpRemove  = "remove"
	UserOpReplace = "replace"
)

// UpdateUser updates and returns the updated user. The API token must have the "Users > Write" scope.
func (c *Client) UpdateUser(id string, operations []UpdateUserOperations) (*User, error) {
	if len(operations) != 0 {
		for _, op := range operations {
			if op.Op != UserOpAdd && op.Op != UserOpRemove && op.Op != UserOpReplace {
				return nil, errors.New("invalid operation: must be 'add', 'remove', or 'replace'")
			}
		}
	} else {
		return nil, errors.New("no operations provided")
	}

	type body struct {
		Operations []UpdateUserOperations `json:"operations"`
	}

	var requestBody body
	requestBody.Operations = operations

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/users/%s", c.BaseURL, id), bytes.NewBuffer(requestBodyJSON))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	fmt.Println(resp)

	return nil, nil
}
