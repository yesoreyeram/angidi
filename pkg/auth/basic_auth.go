package auth

import (
	"golang.org/x/crypto/bcrypt"
	
	"github.com/yesoreyeram/angidi/internal/domain"
	"github.com/yesoreyeram/angidi/internal/repository"
)

// BasicAuthCredentials represents basic authentication credentials
type BasicAuthCredentials struct {
	Email    string
	Password string
}

// BasicAuthStrategy implements basic username/password authentication
type BasicAuthStrategy struct {
	userRepo repository.UserRepository
}

// NewBasicAuthStrategy creates a new basic authentication strategy
func NewBasicAuthStrategy(userRepo repository.UserRepository) *BasicAuthStrategy {
	return &BasicAuthStrategy{
		userRepo: userRepo,
	}
}

// Name returns the name of this authentication strategy
func (s *BasicAuthStrategy) Name() string {
	return "basic"
}

// Authenticate authenticates a user using email and password
func (s *BasicAuthStrategy) Authenticate(credentials interface{}) (*domain.User, error) {
	creds, ok := credentials.(*BasicAuthCredentials)
	if !ok {
		return nil, ErrInvalidCredentials
	}
	
	// Find user by email
	user, err := s.userRepo.FindByEmail(creds.Email)
	if err != nil {
		return nil, ErrInvalidCredentials // Don't reveal if user exists
	}
	
	// Check if account is locked
	if user.IsLocked() {
		return nil, ErrAccountLocked
	}
	
	// Check account status
	if user.Status == domain.UserStatusDisabled {
		return nil, ErrAccountDisabled
	}
	
	if user.Status == domain.UserStatusPendingVerification {
		return nil, ErrAccountPendingVerification
	}
	
	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password))
	if err != nil {
		// Increment failed login count
		s.userRepo.IncrementFailedLogin(user.ID)
		return nil, ErrInvalidCredentials
	}
	
	// Reset failed login count on successful authentication
	s.userRepo.ResetFailedLogin(user.ID)
	
	return user, nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ValidatePasswordStrength validates that a password meets minimum requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	
	// Check for at least one uppercase, one lowercase, and one number
	hasUpper := false
	hasLower := false
	hasNumber := false
	
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}
	
	if !hasUpper || !hasLower || !hasNumber {
		return ErrWeakPassword
	}
	
	return nil
}
