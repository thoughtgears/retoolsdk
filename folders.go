package retoolsdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type Folder struct {
	ID             string `json:"id"`
	LegacyID       string `json:"legacy_id"`
	Name           string `json:"name"`
	ParentFolderID string `json:"parent_folder_id"`
	IsSystemFolder bool   `json:"is_system_folder"`
	FolderType     string `json:"folder_type"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// Folder types allowed.
const (
	FolderTypeWorkflow = "workflow"
	FolderTypeApp      = "app"
	FolderTypeResource = "resource"
)

type FolderType string

func (u *FolderType) String() string {
	return string(*u)
}

// Validate ensures that the options provided in FolderType have valid values.
func (u *FolderType) Validate() error {
	validTypes := map[string]struct{}{
		FolderTypeWorkflow: {},
		FolderTypeApp:      {},
		FolderTypeResource: {},
	}

	// Validate FolderType
	if _, ok := validTypes[string(*u)]; !ok && string(*u) != "" {
		return fmt.Errorf("invalid value for UserType: %s", string(*u))
	}

	return nil
}

// GetFolder returns the folder with the given ID. The API token must have the "Folders > Read" scope.
func (c *Client) GetFolder(id string) (*Folder, error) {
	baseURL := fmt.Sprintf("%s/folders/%s", c.BaseURL, id)
	return doSingleRequest[Folder](c, "GET", baseURL, nil)
}

// ListFolders returns a list of folders. The API token must have the "Folders > Read" scope.
func (c *Client) ListFolders() ([]Folder, error) {
	baseURL := fmt.Sprintf("%s/folders", c.BaseURL)
	return doPaginatedRequest[Folder](c, "GET", baseURL, nil, url.Values{})
}

// CreateFolder creates and returns a folder. The API token must have the "Folders > Write" scope.
func (c *Client) CreateFolder(name, parentFolderID string, folderType FolderType) (*Folder, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if err := folderType.Validate(); err != nil {
		return nil, fmt.Errorf("validating folder type: %w", err)
	}

	requestBody := struct {
		Name           string `json:"name"`
		ParentFolderID string `json:"parent_folder_id,omitempty"`
		FolderType     string `json:"folder_type"`
	}{
		Name:       name,
		FolderType: folderType.String(),
	}

	if parentFolderID != "" {
		requestBody.ParentFolderID = parentFolderID
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	baseURL := fmt.Sprintf("%s/folders", c.BaseURL)
	return doSingleRequest[Folder](c, "POST", baseURL, requestBodyJSON)
}

// UpdateFolder updates a folder by ID. The API token must have the "Folders > Write" scope.
func (c *Client) UpdateFolder(id string, operations []UpdateOperations) (*Folder, error) {
	if len(operations) == 0 {
		return nil, errors.New("no operations provided")
	}

	for _, op := range operations {
		if err := op.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed for operation: %w", err)
		}
	}

	requestBody := struct {
		Operations []UpdateOperations `json:"operations"`
	}{
		Operations: operations,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	baseURL := fmt.Sprintf("%s/folders/%s", c.BaseURL, id)
	return doSingleRequest[Folder](c, "PATCH", baseURL, requestBodyJSON)
}

// DeleteFolder deletes a folder by ID. The API token must have the "Folders > Write" scope.
func (c *Client) DeleteFolder(id string) error {
	baseURL := fmt.Sprintf("%s/folders/%s", c.BaseURL, id)
	_, err := doSingleRequest[any](c, "DELETE", baseURL, nil)
	return err
}
