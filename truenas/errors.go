package truenas

import "fmt"

// NotFoundError represents an error when a resource is not found
type NotFoundError struct {
	ResourceType string
	Identifier   string
}

// Error implements the error interface
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with %s not found", e.ResourceType, e.Identifier)
}

// Is implements error matching for errors.Is()
func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(resourceType, identifier string) *NotFoundError {
	return &NotFoundError{
		ResourceType: resourceType,
		Identifier:   identifier,
	}
}
