package auth

import "github.com/yesoreyeram/angidi/internal/domain"

// AuthStrategy defines the interface for different authentication mechanisms
type AuthStrategy interface {
	// Authenticate verifies the credentials and returns the user if successful
	Authenticate(credentials interface{}) (*domain.User, error)
	// Name returns the name of the authentication strategy
	Name() string
}

// AuthenticationManager manages different authentication strategies
type AuthenticationManager struct {
	strategies map[string]AuthStrategy
	defaultStrategy string
}

// NewAuthenticationManager creates a new authentication manager
func NewAuthenticationManager() *AuthenticationManager {
	return &AuthenticationManager{
		strategies: make(map[string]AuthStrategy),
	}
}

// RegisterStrategy registers a new authentication strategy
func (am *AuthenticationManager) RegisterStrategy(strategy AuthStrategy) {
	am.strategies[strategy.Name()] = strategy
	if am.defaultStrategy == "" {
		am.defaultStrategy = strategy.Name()
	}
}

// SetDefaultStrategy sets the default authentication strategy
func (am *AuthenticationManager) SetDefaultStrategy(name string) error {
	if _, exists := am.strategies[name]; !exists {
		return ErrStrategyNotFound
	}
	am.defaultStrategy = name
	return nil
}

// Authenticate authenticates using the specified strategy (or default if not specified)
func (am *AuthenticationManager) Authenticate(strategyName string, credentials interface{}) (*domain.User, error) {
	if strategyName == "" {
		strategyName = am.defaultStrategy
	}
	
	strategy, exists := am.strategies[strategyName]
	if !exists {
		return nil, ErrStrategyNotFound
	}
	
	return strategy.Authenticate(credentials)
}

// GetStrategy returns a strategy by name
func (am *AuthenticationManager) GetStrategy(name string) (AuthStrategy, error) {
	strategy, exists := am.strategies[name]
	if !exists {
		return nil, ErrStrategyNotFound
	}
	return strategy, nil
}
