package retoolsdk

import (
	"encoding/json"
	"errors"
	"fmt"
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
	baseURL := fmt.Sprintf("%s/users/%s", c.BaseURL, id)
	return doSingleRequest[User](c, "GET", baseURL, nil)
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

	return doPaginatedRequest[User](c, baseURL, query)
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

// Validate ensures that the options provided in CreateUserOpts have valid values.
func (u *User) Validate() error {
	validUserTypes := map[string]struct{}{
		UserTypeDefault: {},
		UserTypeEmbed:   {},
		UserTypeMobile:  {},
	}

	// Validate UserType
	if _, ok := validUserTypes[u.UserType]; !ok && u.UserType != "" {
		return fmt.Errorf("invalid value for UserType: %s", u.UserType)
	}

	return nil
}

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
		} else {
			newUser.UserType = opts.Type
		}
		newUser.Active = opts.Active
	} else {
		newUser.UserType = UserTypeDefault
		newUser.Active = false
	}

	if err := newUser.Validate(); err != nil {
		return nil, err
	}

	baseURL := fmt.Sprintf("%s/users", c.BaseURL)
	return doSingleRequest[User](c, "POST", baseURL, newUser)
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

	baseURL := fmt.Sprintf("%s/users/%s", c.BaseURL, id)
	return doSingleRequest[User](c, "PATCH", baseURL, requestBodyJSON)
}

// DeleteUser disables a user from the organization. The API token must have the "Users > Write" scope.
func (c *Client) DeleteUser(id string) error {
	baseURL := fmt.Sprintf("%s/users/%s", c.BaseURL, id)
	_, err := doSingleRequest[any](c, "DELETE", baseURL, nil)
	return err
}
