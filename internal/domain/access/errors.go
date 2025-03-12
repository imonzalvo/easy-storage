package access

import "errors"

// Add this error type if it doesn't already exist
var (
	ErrInvalidResourceType   = errors.New("invalid resource type")
	ErrServiceNotInitialized = errors.New("service not properly initialized")
)
