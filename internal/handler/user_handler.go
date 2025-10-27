package handler

import (
	"encoding/json"
	"net/http"
	
	"github.com/yesoreyeram/angidi/internal/config"
	"github.com/yesoreyeram/angidi/internal/domain"
	"github.com/yesoreyeram/angidi/internal/service"
	"github.com/yesoreyeram/angidi/pkg/auth"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService *service.UserService
	config      *config.Config
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		config:      cfg,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, domain.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}
	
	user, err := h.userService.Register(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to register user"
		
		if err == auth.ErrUserAlreadyExists {
			statusCode = http.StatusConflict
			message = "User with this email already exists"
		} else if err == auth.ErrWeakPassword {
			statusCode = http.StatusBadRequest
			message = "Password must be at least 8 characters with uppercase, lowercase, and number"
		}
		
		respondJSON(w, statusCode, domain.ErrorResponse{
			Error:   "registration_failed",
			Message: message,
		})
		return
	}
	
	respondJSON(w, http.StatusCreated, domain.SuccessResponse{
		Message: "User registered successfully",
		Data:    user,
	})
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, domain.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}
	
	user, err := h.userService.Login(&req, "basic")
	if err != nil {
		statusCode := http.StatusUnauthorized
		message := "Invalid credentials"
		
		if err == auth.ErrAccountLocked {
			statusCode = http.StatusForbidden
			message = "Account is locked due to too many failed login attempts. Please try again later."
		} else if err == auth.ErrAccountDisabled {
			statusCode = http.StatusForbidden
			message = "Account is disabled"
		} else if err == auth.ErrAccountPendingVerification {
			statusCode = http.StatusForbidden
			message = "Account is pending verification"
		}
		
		respondJSON(w, statusCode, domain.ErrorResponse{
			Error:   "login_failed",
			Message: message,
		})
		return
	}
	
	// Set session cookie
	if req.RememberMe {
		h.setSessionCookie(w, user.ID, h.config.Session.CookieMaxAge)
	} else {
		// Session expires when browser closes
		h.setSessionCookie(w, user.ID, 0)
	}
	
	respondJSON(w, http.StatusOK, domain.LoginResponse{
		User:    user,
		Message: "Login successful",
	})
}

// Logout handles user logout
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     h.config.Session.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: h.config.Session.CookieHttpOnly,
		Secure:   h.config.Session.CookieSecure,
	})
	
	respondJSON(w, http.StatusOK, domain.SuccessResponse{
		Message: "Logout successful",
	})
}

// ForgotPassword handles forgot password requests
func (h *UserHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req domain.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, domain.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}
	
	_, err := h.userService.ForgotPassword(&req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
			Error:   "password_reset_failed",
			Message: "Failed to process password reset request",
		})
		return
	}
	
	// Always return success to prevent email enumeration
	// In production, the token would be sent via email
	respondJSON(w, http.StatusOK, domain.SuccessResponse{
		Message: "If an account exists with this email, password reset instructions have been sent",
	})
}

// ResetPassword handles password reset with token
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req domain.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, domain.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}
	
	err := h.userService.ResetPassword(&req)
	if err != nil {
		statusCode := http.StatusBadRequest
		message := "Failed to reset password"
		
		if err == auth.ErrInvalidToken {
			message = "Invalid or expired reset token"
		} else if err == auth.ErrWeakPassword {
			message = "Password must be at least 8 characters with uppercase, lowercase, and number"
		}
		
		respondJSON(w, statusCode, domain.ErrorResponse{
			Error:   "password_reset_failed",
			Message: message,
		})
		return
	}
	
	respondJSON(w, http.StatusOK, domain.SuccessResponse{
		Message: "Password reset successful",
	})
}

// setSessionCookie sets a session cookie
// NOTE: This is a simplified session implementation for demonstration purposes.
// In production, use:
// 1. Signed/encrypted session tokens (e.g., using gorilla/sessions)
// 2. Random session IDs mapped to user data in a session store (Redis, database)
// 3. Session rotation on privilege escalation
// 4. Proper session expiration and cleanup
func (h *UserHandler) setSessionCookie(w http.ResponseWriter, userID string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.config.Session.CookieName,
		Value:    userID,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: h.config.Session.CookieHttpOnly,
		Secure:   h.config.Session.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

// respondJSON writes a JSON response
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// GetCurrentUser returns the current logged-in user
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.config.Session.CookieName)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, domain.ErrorResponse{
			Error:   "not_authenticated",
			Message: "Not authenticated",
		})
		return
	}
	
	user, err := h.userService.GetUserByID(cookie.Value)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, domain.ErrorResponse{
			Error:   "not_authenticated",
			Message: "Invalid session",
		})
		return
	}
	
	respondJSON(w, http.StatusOK, user)
}
