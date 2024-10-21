package retoolsdk

import (
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
)

// Subject represents the subject in the response (group, user, userInvite).
type Subject struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	AccessLevel string `json:"access_level,omitempty"`
}

// Sources represents the sources of access (direct, universal, groups, inherited).
type Sources struct {
	Direct    bool      `json:"direct"`
	Universal bool      `json:"universal"`
	Groups    []Subject `json:"groups"`
	Inherited Subject   `json:"inherited"`
}

// AccessData represents the access information for group, user, or userInvite.
type AccessData struct {
	Subject     Subject `json:"subject"`
	Sources     Sources `json:"sources"`
	AccessLevel string  `json:"accessLevel"`
}

// GroupedData holds the different categories of access data returned in the "data" field.
type GroupedData struct {
	Group      []AccessData `json:"group"`
	User       []AccessData `json:"user"`
	UserInvite []AccessData `json:"userInvite"`
}

// Access levels allowed.
const (
	NoneAccess = "none"
	UseAccess  = "use"
	EditAccess = "edit"
	OwnAccess  = "own"
)

type AccessLevel string

// Validate ensures that the provided access level has a valid value.
func (a *AccessLevel) Validate() error {
	validAccessLevels := map[string]struct{}{
		NoneAccess: {},
		UseAccess:  {},
		EditAccess: {},
		OwnAccess:  {},
	}

	if _, ok := validAccessLevels[string(*a)]; !ok {
		return fmt.Errorf("invalid access level: %s", string(*a))
	}

	return nil
}

func (a *AccessLevel) String() string {
	return string(*a)
}

// Object types allowed.
const (
	AppObject                   = "app"
	FolderObject                = "folder"
	ResourceObject              = "resource"
	ResourceConfigurationObject = "resourceConfiguration"
)

type ObjectType string

func (o *ObjectType) String() string {
	return string(*o)
}

// Validate ensures that the provided object types has a valid value.
func (o *ObjectType) Validate() error {
	validTypes := map[string]struct{}{
		AppObject:                   {},
		FolderObject:                {},
		ResourceObject:              {},
		ResourceConfigurationObject: {},
	}

	if _, ok := validTypes[string(*o)]; !ok {
		return fmt.Errorf("invalid object type: %s", string(*o))
	}

	return nil
}

// GetFolderOrAppAccessList Returns the list of users/groups and corresponding access levels whom have access to a
// selected folder/page. The API token must have the "Permissions > Read" scope.
// Supported from onprem edge version 3.96.0+ and 3.114-stable+.
func (c *Client) GetFolderOrAppAccessList(objectID, objectType ObjectType) (*GroupedData, error) {
	if err := objectType.Validate(); err != nil {
		return nil, fmt.Errorf("validating object type: %w", err)
	}

	baseURL := fmt.Sprintf("%s/permissions/accessList/%s/%s", c.BaseURL, objectType, objectID)
	return doSingleRequest[GroupedData](c, "GET", baseURL, nil)
}

// ListGroupObjectPermissions returns the list of objects with corresponding access levels that a subject (group) has
// access to. The API token must have the "Permissions > Read" scope.
// Folders are supported from API version 2.0.0 + and onprem version 3.18+,
// apps are supported from API version 2.4.0+ and onprem version 3.26.0+,
// resources and resource_configurations are supported from onprem edge version 3.37.0+ and 3.47-stable+.
func (c *Client) ListGroupObjectPermissions(subject string, objectType ObjectType, id any) ([]Subject, error) {
	if err := objectType.Validate(); err != nil {
		return nil, fmt.Errorf("validating object type: %w", err)
	}

	if !slices.Contains([]string{"group", "user"}, subject) {
		return nil, fmt.Errorf("invalid subject: %s", subject)
	}

	var requestBody any

	if subject == "group" {
		groupID, ok := id.(int)
		if !ok {
			return nil, fmt.Errorf("invalid id type for group: expected int")
		}

		requestBody = struct {
			Subject struct {
				ID   int    `json:"id"`
				Type string `json:"type"`
			}
			ObjectType ObjectType `json:"objectType"`
		}{
			Subject: struct {
				ID   int    `json:"id"`
				Type string `json:"type"`
			}{
				ID:   groupID,
				Type: subject,
			},
			ObjectType: objectType,
		}
	} else if subject == "user" {
		userID, ok := id.(string)
		if !ok {
			return nil, fmt.Errorf("invalid id type for user: expected string")
		}

		requestBody = struct {
			Subject struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}
			ObjectType ObjectType `json:"objectType"`
		}{
			Subject: struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}{
				ID:   userID,
				Type: subject,
			},
			ObjectType: objectType,
		}
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	baseURL := fmt.Sprintf("%s/permissions/listObjects", c.BaseURL)
	return doPaginatedRequest[Subject](c, "POST", baseURL, requestBodyJSON, url.Values{})
}

func (c *Client) GrantPermission(subject string, subjectID any, objectType ObjectType, objectID string, accessLevel AccessLevel) ([]Subject, error) {
	if err := objectType.Validate(); err != nil {
		return nil, fmt.Errorf("validating object type: %w", err)
	}

	if !slices.Contains([]string{"group", "user"}, subject) {
		return nil, fmt.Errorf("invalid subject: %s", subject)
	}

	var requestBody any

	if subject == "group" {
		id, ok := subjectID.(int)
		if !ok {
			return nil, fmt.Errorf("invalid id type for group: expected int")
		}

		requestBody = struct {
			Subject struct {
				ID   int    `json:"id"`
				Type string `json:"type"`
			}
			Object struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}
			AccessLevel string `json:"access_level"`
		}{
			Subject: struct {
				ID   int    `json:"id"`
				Type string `json:"type"`
			}{
				ID:   id,
				Type: subject,
			},
			Object: struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}{
				ID:   objectID,
				Type: objectType.String(),
			},
			AccessLevel: accessLevel.String(),
		}
	} else if subject == "user" {
		id, ok := subjectID.(string)
		if !ok {
			return nil, fmt.Errorf("invalid id type for user: expected string")
		}

		requestBody = struct {
			Subject struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}
			Object struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}
			AccessLevel string `json:"access_level"`
		}{
			Subject: struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}{
				ID:   id,
				Type: subject,
			},
			Object: struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}{
				ID:   objectID,
				Type: objectType.String(),
			},
			AccessLevel: accessLevel.String(),
		}
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	baseURL := fmt.Sprintf("%s/permissions/grant", c.BaseURL)
	return doPaginatedRequest[Subject](c, "POST", baseURL, requestBodyJSON, url.Values{})
}

func (c *Client) RevokePermission(subject string, subjectID any, objectType ObjectType, objectID string) ([]Subject, error) {
	if err := objectType.Validate(); err != nil {
		return nil, fmt.Errorf("validating object type: %w", err)
	}

	if !slices.Contains([]string{"group", "user"}, subject) {
		return nil, fmt.Errorf("invalid subject: %s", subject)
	}

	var requestBody any

	if subject == "group" {
		id, ok := subjectID.(int)
		if !ok {
			return nil, fmt.Errorf("invalid id type for group: expected int")
		}

		requestBody = struct {
			Subject struct {
				ID   int    `json:"id"`
				Type string `json:"type"`
			}
			Object struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}
		}{
			Subject: struct {
				ID   int    `json:"id"`
				Type string `json:"type"`
			}{
				ID:   id,
				Type: subject,
			},
			Object: struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}{
				ID:   objectID,
				Type: objectType.String(),
			},
		}
	} else if subject == "user" {
		id, ok := subjectID.(string)
		if !ok {
			return nil, fmt.Errorf("invalid id type for user: expected string")
		}

		requestBody = struct {
			Subject struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}
			Object struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}
		}{
			Subject: struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}{
				ID:   id,
				Type: subject,
			},
			Object: struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			}{
				ID:   objectID,
				Type: objectType.String(),
			},
		}
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}
	baseURL := fmt.Sprintf("%s/permissions/revoke", c.BaseURL)
	return doPaginatedRequest[Subject](c, "POST", baseURL, requestBodyJSON, url.Values{})
}
