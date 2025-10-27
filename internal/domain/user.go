package domain

import "time"

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusActive             UserStatus = "active"
	UserStatusPendingVerification UserStatus = "pending_verification"
	UserStatusDisabled           UserStatus = "disabled"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleSuperUser UserRole = "superuser"
	UserRoleAdmin     UserRole = "admin"
	UserRoleViewer    UserRole = "viewer"
	UserRoleCustomer  UserRole = "customer"
)

// User represents a user in the system
type User struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"` // Never expose password hash in JSON
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Role              UserRole   `json:"role"`
	Status            UserStatus `json:"status"`
	FailedLoginCount  int        `json:"-"`
	LastFailedLogin   *time.Time `json:"-"`
	PasswordResetToken string    `json:"-"`
	PasswordResetExpiry *time.Time `json:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsLocked checks if the user account is locked due to too many failed login attempts
func (u *User) IsLocked() bool {
	if u.FailedLoginCount >= 5 && u.LastFailedLogin != nil {
		// Lock for 30 minutes after 5 failed attempts
		lockDuration := 30 * time.Minute
		return time.Since(*u.LastFailedLogin) < lockDuration
	}
	return false
}

// CanResetPassword checks if the password reset token is valid
func (u *User) CanResetPassword() bool {
	return u.PasswordResetToken != "" && u.PasswordResetExpiry != nil && time.Now().Before(*u.PasswordResetExpiry)
}
