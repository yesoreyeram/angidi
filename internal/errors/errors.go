package errors

import "errors"

var (
	// ErrStrategyNotFound is returned when the requested authentication strategy is not found
	ErrStrategyNotFound = errors.New("authentication strategy not found")
	
	// ErrInvalidCredentials is returned when the provided credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	
	// ErrAccountLocked is returned when the account is locked due to too many failed attempts
	ErrAccountLocked = errors.New("account is locked due to too many failed login attempts")
	
	// ErrAccountDisabled is returned when the account is disabled
	ErrAccountDisabled = errors.New("account is disabled")
	
	// ErrAccountPendingVerification is returned when the account is pending verification
	ErrAccountPendingVerification = errors.New("account is pending verification")
	
	// ErrUserNotFound is returned when the user is not found
	ErrUserNotFound = errors.New("user not found")
	
	// ErrUserAlreadyExists is returned when a user with the same email already exists
	ErrUserAlreadyExists = errors.New("user already exists")
	
	// ErrInvalidToken is returned when the provided token is invalid
	ErrInvalidToken = errors.New("invalid or expired token")
	
	// ErrWeakPassword is returned when the password doesn't meet strength requirements
	ErrWeakPassword = errors.New("password does not meet strength requirements")
)
