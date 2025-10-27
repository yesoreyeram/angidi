package repository

import (
	"sync"
	"time"
	
	"github.com/yesoreyeram/angidi/internal/domain"
	"github.com/yesoreyeram/angidi/internal/errors"
)

// InMemoryUserRepository implements UserRepository using in-memory storage
type InMemoryUserRepository struct {
	users      map[string]*domain.User
	emailIndex map[string]string // email -> user ID
	tokenIndex map[string]string // reset token -> user ID
	mu         sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:      make(map[string]*domain.User),
		emailIndex: make(map[string]string),
		tokenIndex: make(map[string]string),
	}
}

// Create creates a new user
func (r *InMemoryUserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if user with same email already exists
	if _, exists := r.emailIndex[user.Email]; exists {
		return errors.ErrUserAlreadyExists
	}
	
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	r.users[user.ID] = user
	r.emailIndex[user.Email] = user.ID
	
	return nil
}

// FindByID finds a user by ID
func (r *InMemoryUserRepository) FindByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, errors.ErrUserNotFound
	}
	
	return user, nil
}

// FindByEmail finds a user by email
func (r *InMemoryUserRepository) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	userID, exists := r.emailIndex[email]
	if !exists {
		return nil, errors.ErrUserNotFound
	}
	
	return r.users[userID], nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	existingUser, exists := r.users[user.ID]
	if !exists {
		return errors.ErrUserNotFound
	}
	
	// Update email index if email changed
	if existingUser.Email != user.Email {
		delete(r.emailIndex, existingUser.Email)
		r.emailIndex[user.Email] = user.ID
	}
	
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	
	return nil
}

// Delete deletes a user by ID
func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	user, exists := r.users[id]
	if !exists {
		return errors.ErrUserNotFound
	}
	
	delete(r.emailIndex, user.Email)
	if user.PasswordResetToken != "" {
		delete(r.tokenIndex, user.PasswordResetToken)
	}
	delete(r.users, id)
	
	return nil
}

// List returns all users
func (r *InMemoryUserRepository) List() ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	
	return users, nil
}

// IncrementFailedLogin increments the failed login counter for a user
func (r *InMemoryUserRepository) IncrementFailedLogin(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	user, exists := r.users[id]
	if !exists {
		return errors.ErrUserNotFound
	}
	
	user.FailedLoginCount++
	now := time.Now()
	user.LastFailedLogin = &now
	user.UpdatedAt = time.Now()
	
	return nil
}

// ResetFailedLogin resets the failed login counter for a user
func (r *InMemoryUserRepository) ResetFailedLogin(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	user, exists := r.users[id]
	if !exists {
		return errors.ErrUserNotFound
	}
	
	user.FailedLoginCount = 0
	user.LastFailedLogin = nil
	user.UpdatedAt = time.Now()
	
	return nil
}

// SetPasswordResetToken sets a password reset token for a user
func (r *InMemoryUserRepository) SetPasswordResetToken(email string, token string, expiry int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	userID, exists := r.emailIndex[email]
	if !exists {
		return errors.ErrUserNotFound
	}
	
	user := r.users[userID]
	
	// Remove old token from index if exists
	if user.PasswordResetToken != "" {
		delete(r.tokenIndex, user.PasswordResetToken)
	}
	
	expiryTime := time.Unix(expiry, 0)
	user.PasswordResetToken = token
	user.PasswordResetExpiry = &expiryTime
	user.UpdatedAt = time.Now()
	
	r.tokenIndex[token] = user.ID
	
	return nil
}

// FindByPasswordResetToken finds a user by password reset token
func (r *InMemoryUserRepository) FindByPasswordResetToken(token string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	userID, exists := r.tokenIndex[token]
	if !exists {
		return nil, errors.ErrInvalidToken
	}
	
	return r.users[userID], nil
}

// ClearPasswordResetToken clears the password reset token for a user
func (r *InMemoryUserRepository) ClearPasswordResetToken(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	user, exists := r.users[id]
	if !exists {
		return errors.ErrUserNotFound
	}
	
	if user.PasswordResetToken != "" {
		delete(r.tokenIndex, user.PasswordResetToken)
	}
	
	user.PasswordResetToken = ""
	user.PasswordResetExpiry = nil
	user.UpdatedAt = time.Now()
	
	return nil
}
