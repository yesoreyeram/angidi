package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	
	"github.com/google/uuid"
	
	"github.com/yesoreyeram/angidi/internal/domain"
	"github.com/yesoreyeram/angidi/internal/errors"
	"github.com/yesoreyeram/angidi/internal/repository"
	"github.com/yesoreyeram/angidi/pkg/auth"
)

// UserService provides user-related business logic
type UserService struct {
	userRepo repository.UserRepository
	authMgr  *auth.AuthenticationManager
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, authMgr *auth.AuthenticationManager) *UserService {
	return &UserService{
		userRepo: userRepo,
		authMgr:  authMgr,
	}
}

// Register registers a new user with the default role of Customer/Viewer
func (s *UserService) Register(req *domain.RegisterRequest) (*domain.User, error) {
	// Validate password strength
	if err := auth.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}
	
	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Create user
	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         domain.UserRoleCustomer, // Default role
		Status:       domain.UserStatusActive,  // Could be PendingVerification if email verification is required
	}
	
	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	
	return user, nil
}

// Login authenticates a user using the specified strategy
func (s *UserService) Login(req *domain.LoginRequest, strategyName string) (*domain.User, error) {
	if strategyName == "" {
		strategyName = "basic"
	}
	
	credentials := &auth.BasicAuthCredentials{
		Email:    req.Email,
		Password: req.Password,
	}
	
	user, err := s.authMgr.Authenticate(strategyName, credentials)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// ForgotPassword initiates the password reset process
func (s *UserService) ForgotPassword(req *domain.ForgotPasswordRequest) (string, error) {
	// Generate a random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)
	
	// Set token expiry to 1 hour from now
	expiry := time.Now().Add(1 * time.Hour).Unix()
	
	// Save token to user
	err := s.userRepo.SetPasswordResetToken(req.Email, token, expiry)
	if err != nil {
		// Don't reveal if user exists or not
		if err == errors.ErrUserNotFound {
			return "", nil
		}
		return "", err
	}
	
	// In a real application, send email with reset link here
	// For now, we'll just return the token
	return token, nil
}

// ResetPassword resets a user's password using a valid reset token
func (s *UserService) ResetPassword(req *domain.ResetPasswordRequest) error {
	// Validate new password strength
	if err := auth.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}
	
	// Find user by token
	user, err := s.userRepo.FindByPasswordResetToken(req.Token)
	if err != nil {
		return errors.ErrInvalidToken
	}
	
	// Check if token is still valid
	if !user.CanResetPassword() {
		return errors.ErrInvalidToken
	}
	
	// Hash new password
	passwordHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Update password
	user.PasswordHash = passwordHash
	if err := s.userRepo.Update(user); err != nil {
		return err
	}
	
	// Clear reset token
	if err := s.userRepo.ClearPasswordResetToken(user.ID); err != nil {
		return err
	}
	
	return nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id string) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(email)
}

// CreateSuperUser creates the initial superuser account
func (s *UserService) CreateSuperUser(email, password, firstName, lastName string) (*domain.User, error) {
	// Check if superuser already exists
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		// Superuser already exists
		return nil, fmt.Errorf("superuser already exists")
	}
	
	// Hash password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Create superuser
	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Role:         domain.UserRoleSuperUser,
		Status:       domain.UserStatusActive,
	}
	
	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	
	return user, nil
}
