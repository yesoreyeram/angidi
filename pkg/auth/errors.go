package auth

// Re-export errors from internal/errors package for backward compatibility
import apperrors "github.com/yesoreyeram/angidi/internal/errors"

var (
	ErrStrategyNotFound           = apperrors.ErrStrategyNotFound
	ErrInvalidCredentials         = apperrors.ErrInvalidCredentials
	ErrAccountLocked              = apperrors.ErrAccountLocked
	ErrAccountDisabled            = apperrors.ErrAccountDisabled
	ErrAccountPendingVerification = apperrors.ErrAccountPendingVerification
	ErrUserNotFound               = apperrors.ErrUserNotFound
	ErrUserAlreadyExists          = apperrors.ErrUserAlreadyExists
	ErrInvalidToken               = apperrors.ErrInvalidToken
	ErrWeakPassword               = apperrors.ErrWeakPassword
)
