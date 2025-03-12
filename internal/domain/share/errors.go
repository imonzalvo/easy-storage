package share

import "errors"

var (
	// ErrShareNotFound is returned when a share is not found
	ErrShareNotFound = errors.New("share not found")

	// ErrShareRevoked is returned when attempting to access a revoked share
	ErrShareRevoked = errors.New("share has been revoked")

	// ErrShareExpired is returned when attempting to access an expired share
	ErrShareExpired = errors.New("share has expired")

	// ErrInvalidPassword is returned when an incorrect password is provided
	ErrInvalidPassword = errors.New("invalid share password")

	// ErrInvalidShareType is returned when an operation is attempted on the wrong share type
	ErrInvalidShareType = errors.New("invalid share type for this operation")

	// ErrUnauthorizedAccess is returned when a user attempts to access a share they don't have permission for
	ErrUnauthorizedAccess = errors.New("unauthorized access to share")
)
