# JWT Authentication Implementation Guide

This guide explains how to add JWT (JSON Web Token) authentication to the existing authentication system using the Strategy Pattern.

## Overview

JWT authentication will be implemented as a new authentication strategy alongside the existing Basic Auth strategy, demonstrating the extensibility of the Strategy Pattern design.

## Dependencies

Add JWT library to your project:

```bash
go get github.com/golang-jwt/jwt/v5
```

## Implementation Steps

### 1. Create JWT Configuration

Update `internal/config/config.go`:

```go
type JWTConfig struct {
    SecretKey            string
    AccessTokenExpiry    int // in minutes
    RefreshTokenExpiry   int // in hours
    Issuer               string
}

func Load() *Config {
    return &Config{
        // ... existing config
        JWT: JWTConfig{
            SecretKey:          getEnv("JWT_SECRET", "your-256-bit-secret"),
            AccessTokenExpiry:  getEnvAsInt("JWT_ACCESS_EXPIRY", 15),  // 15 minutes
            RefreshTokenExpiry: getEnvAsInt("JWT_REFRESH_EXPIRY", 168), // 7 days
            Issuer:             getEnv("JWT_ISSUER", "angidi"),
        },
    }
}
```

### 2. Create JWT Service

Create `pkg/auth/jwt_service.go`:

```go
package auth

import (
    "fmt"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/yesoreyeram/angidi/internal/domain"
)

type JWTService struct {
    secretKey          []byte
    accessTokenExpiry  time.Duration
    refreshTokenExpiry time.Duration
    issuer             string
}

type CustomClaims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func NewJWTService(secretKey string, accessExpiry, refreshExpiry int, issuer string) *JWTService {
    return &JWTService{
        secretKey:          []byte(secretKey),
        accessTokenExpiry:  time.Duration(accessExpiry) * time.Minute,
        refreshTokenExpiry: time.Duration(refreshExpiry) * time.Hour,
        issuer:             issuer,
    }
}

func (s *JWTService) GenerateAccessToken(user *domain.User) (string, error) {
    claims := CustomClaims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   string(user.Role),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    s.issuer,
            Subject:   user.ID,
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *JWTService) GenerateRefreshToken(user *domain.User) (string, error) {
    claims := jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenExpiry)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        NotBefore: jwt.NewNumericDate(time.Now()),
        Issuer:    s.issuer,
        Subject:   user.ID,
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}
```

### 3. Create JWT Authentication Strategy

Create `pkg/auth/jwt_auth.go`:

```go
package auth

import (
    "strings"
    
    "github.com/yesoreyeram/angidi/internal/domain"
    "github.com/yesoreyeram/angidi/internal/repository"
)

type JWTAuthCredentials struct {
    Token string
}

type JWTAuthStrategy struct {
    userRepo   repository.UserRepository
    jwtService *JWTService
}

func NewJWTAuthStrategy(userRepo repository.UserRepository, jwtService *JWTService) *JWTAuthStrategy {
    return &JWTAuthStrategy{
        userRepo:   userRepo,
        jwtService: jwtService,
    }
}

func (s *JWTAuthStrategy) Name() string {
    return "jwt"
}

func (s *JWTAuthStrategy) Authenticate(credentials interface{}) (*domain.User, error) {
    creds, ok := credentials.(*JWTAuthCredentials)
    if !ok {
        return nil, ErrInvalidCredentials
    }
    
    // Remove "Bearer " prefix if present
    token := strings.TrimPrefix(creds.Token, "Bearer ")
    
    // Validate token
    claims, err := s.jwtService.ValidateToken(token)
    if err != nil {
        return nil, ErrInvalidCredentials
    }
    
    // Get user from repository
    user, err := s.userRepo.FindByID(claims.UserID)
    if err != nil {
        return nil, ErrUserNotFound
    }
    
    // Check if account is active
    if user.Status != domain.UserStatusActive {
        return nil, ErrAccountDisabled
    }
    
    return user, nil
}
```

### 4. Update Login Handler

Update `internal/handler/user_handler.go`:

```go
type UserHandler struct {
    userService *service.UserService
    jwtService  *auth.JWTService
    config      *config.Config
}

func NewUserHandler(userService *service.UserService, jwtService *auth.JWTService, cfg *config.Config) *UserHandler {
    return &UserHandler{
        userService: userService,
        jwtService:  jwtService,
        config:      cfg,
    }
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    // ... existing validation code
    
    user, err := h.userService.Login(&req, "basic")
    if err != nil {
        // ... error handling
        return
    }
    
    // Generate JWT tokens
    accessToken, err := h.jwtService.GenerateAccessToken(user)
    if err != nil {
        respondJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
            Error:   "token_generation_failed",
            Message: "Failed to generate access token",
        })
        return
    }
    
    refreshToken, err := h.jwtService.GenerateRefreshToken(user)
    if err != nil {
        respondJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
            Error:   "token_generation_failed",
            Message: "Failed to generate refresh token",
        })
        return
    }
    
    // Set refresh token as HTTP-only cookie
    if req.RememberMe {
        http.SetCookie(w, &http.Cookie{
            Name:     "refresh_token",
            Value:    refreshToken,
            Path:     "/",
            MaxAge:   h.config.JWT.RefreshTokenExpiry * 3600,
            HttpOnly: true,
            Secure:   h.config.Session.CookieSecure,
            SameSite: http.SameSiteStrictMode,
        })
    }
    
    respondJSON(w, http.StatusOK, domain.LoginResponse{
        User:    user,
        Token:   accessToken,
        Message: "Login successful",
    })
}
```

### 5. Create JWT Middleware

Create `internal/middleware/jwt_auth.go`:

```go
package middleware

import (
    "context"
    "net/http"
    "strings"
    
    "github.com/yesoreyeram/angidi/pkg/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

func JWTAuth(jwtService *auth.JWTService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Missing authorization header", http.StatusUnauthorized)
                return
            }
            
            token := strings.TrimPrefix(authHeader, "Bearer ")
            
            claims, err := jwtService.ValidateToken(token)
            if err != nil {
                http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
                return
            }
            
            // Add claims to context
            ctx := context.WithValue(r.Context(), UserContextKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### 6. Add Token Refresh Endpoint

Add to `internal/handler/user_handler.go`:

```go
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("refresh_token")
    if err != nil {
        respondJSON(w, http.StatusUnauthorized, domain.ErrorResponse{
            Error:   "invalid_token",
            Message: "Refresh token not found",
        })
        return
    }
    
    claims, err := h.jwtService.ValidateToken(cookie.Value)
    if err != nil {
        respondJSON(w, http.StatusUnauthorized, domain.ErrorResponse{
            Error:   "invalid_token",
            Message: "Invalid or expired refresh token",
        })
        return
    }
    
    user, err := h.userService.GetUserByID(claims.Subject)
    if err != nil {
        respondJSON(w, http.StatusUnauthorized, domain.ErrorResponse{
            Error:   "user_not_found",
            Message: "User not found",
        })
        return
    }
    
    accessToken, err := h.jwtService.GenerateAccessToken(user)
    if err != nil {
        respondJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
            Error:   "token_generation_failed",
            Message: "Failed to generate access token",
        })
        return
    }
    
    respondJSON(w, http.StatusOK, map[string]string{
        "access_token": accessToken,
    })
}
```

### 7. Update Routes

Update `cmd/user-service/main.go`:

```go
func main() {
    // ... existing setup
    
    // Initialize JWT service
    jwtService := auth.NewJWTService(
        cfg.JWT.SecretKey,
        cfg.JWT.AccessTokenExpiry,
        cfg.JWT.RefreshTokenExpiry,
        cfg.JWT.Issuer,
    )
    
    // Register JWT strategy
    jwtAuth := auth.NewJWTAuthStrategy(userRepo, jwtService)
    authMgr.RegisterStrategy(jwtAuth)
    
    // Initialize handlers with JWT service
    userHandler := handler.NewUserHandler(userService, jwtService, cfg)
    
    // JWT middleware
    jwtMiddleware := middleware.JWTAuth(jwtService)
    
    // Routes
    mux.HandleFunc("/api/register", authRateLimiter.Middleware(userHandler.Register))
    mux.HandleFunc("/api/login", authRateLimiter.Middleware(userHandler.Login))
    mux.HandleFunc("/api/logout", userHandler.Logout)
    mux.HandleFunc("/api/refresh", userHandler.RefreshToken)
    
    // Protected routes
    mux.Handle("/api/me", jwtMiddleware(http.HandlerFunc(userHandler.GetCurrentUser)))
    
    // ... rest of setup
}
```

### 8. Frontend Integration

Update JavaScript to use JWT:

```javascript
// Store token in localStorage or sessionStorage
async function login(email, password, rememberMe) {
    const response = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password, remember_me: rememberMe })
    });
    
    if (response.ok) {
        const data = await response.json();
        // Store access token
        localStorage.setItem('access_token', data.token);
        window.location.href = '/dashboard';
    }
}

// Use token in requests
async function getUser() {
    const token = localStorage.getItem('access_token');
    const response = await fetch('/api/me', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    
    if (response.status === 401) {
        // Token expired, try to refresh
        await refreshToken();
    }
    
    return response.json();
}

// Refresh token
async function refreshToken() {
    const response = await fetch('/api/refresh', {
        method: 'POST',
        credentials: 'include' // Send cookies
    });
    
    if (response.ok) {
        const data = await response.json();
        localStorage.setItem('access_token', data.access_token);
    } else {
        // Refresh failed, redirect to login
        window.location.href = '/login';
    }
}
```

## Security Considerations

1. **Secret Key**: Use a strong, randomly generated secret key (at least 256 bits)
2. **Token Expiry**: Keep access tokens short-lived (15-30 minutes)
3. **HTTPS Only**: Always use HTTPS in production
4. **Token Storage**: 
   - Access tokens: localStorage (accessible to XSS but works with CORS)
   - Refresh tokens: HTTP-only cookies (protected from XSS)
5. **Token Revocation**: Implement a token blacklist for logout/revocation
6. **Refresh Token Rotation**: Issue new refresh token on each use

## Testing

```go
func TestJWTService(t *testing.T) {
    jwtService := auth.NewJWTService("test-secret", 15, 168, "test")
    
    user := &domain.User{
        ID:    "123",
        Email: "test@example.com",
        Role:  domain.UserRoleCustomer,
    }
    
    token, err := jwtService.GenerateAccessToken(user)
    if err != nil {
        t.Fatalf("failed to generate token: %v", err)
    }
    
    claims, err := jwtService.ValidateToken(token)
    if err != nil {
        t.Fatalf("failed to validate token: %v", err)
    }
    
    if claims.UserID != user.ID {
        t.Errorf("expected user ID %s, got %s", user.ID, claims.UserID)
    }
}
```
