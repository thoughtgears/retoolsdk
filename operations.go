package retoolsdk

import "fmt"

// UpdateOperations is a struct that contains the operations to update resources.
type UpdateOperations struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// Update operations allowed.
const (
	OpAdd     = "add"
	OpRemove  = "remove"
	OpReplace = "replace"
)

// Validate ensures that the provided operation types in UpdateOperations have valid values.
func (u *UpdateOperations) Validate() error {
	validOperationTypes := map[string]struct{}{
		OpAdd:     {},
		OpRemove:  {},
		OpReplace: {},
	}

	if _, ok := validOperationTypes[u.Op]; !ok {
		return fmt.Errorf("invalid operation type: %s", u.Op)
	}

	if u.Path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if u.Op != OpRemove && u.Value == "" {
		return fmt.Errorf("value cannot be empty for %s operation", u.Op)
	}

	return nil
}
