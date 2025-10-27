package repository

import "github.com/yesoreyeram/angidi/internal/domain"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(user *domain.User) error
	
	// FindByID finds a user by ID
	FindByID(id string) (*domain.User, error)
	
	// FindByEmail finds a user by email
	FindByEmail(email string) (*domain.User, error)
	
	// Update updates an existing user
	Update(user *domain.User) error
	
	// Delete deletes a user by ID
	Delete(id string) error
	
	// List returns all users
	List() ([]*domain.User, error)
	
	// IncrementFailedLogin increments the failed login counter for a user
	IncrementFailedLogin(id string) error
	
	// ResetFailedLogin resets the failed login counter for a user
	ResetFailedLogin(id string) error
	
	// SetPasswordResetToken sets a password reset token for a user
	SetPasswordResetToken(email string, token string, expiry int64) error
	
	// FindByPasswordResetToken finds a user by password reset token
	FindByPasswordResetToken(token string) (*domain.User, error)
	
	// ClearPasswordResetToken clears the password reset token for a user
	ClearPasswordResetToken(id string) error
}
