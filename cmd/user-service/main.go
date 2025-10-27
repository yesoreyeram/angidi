package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	
	"github.com/yesoreyeram/angidi/internal/config"
	"github.com/yesoreyeram/angidi/internal/handler"
	"github.com/yesoreyeram/angidi/internal/middleware"
	"github.com/yesoreyeram/angidi/internal/repository"
	"github.com/yesoreyeram/angidi/internal/service"
	"github.com/yesoreyeram/angidi/pkg/auth"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize repository
	userRepo := repository.NewInMemoryUserRepository()
	
	// Initialize authentication manager
	authMgr := auth.NewAuthenticationManager()
	basicAuth := auth.NewBasicAuthStrategy(userRepo)
	authMgr.RegisterStrategy(basicAuth)
	
	// Initialize service
	userService := service.NewUserService(userRepo, authMgr)
	
	// Bootstrap superuser
	superUser, err := userService.CreateSuperUser(
		cfg.SuperUser.Email,
		cfg.SuperUser.Password,
		cfg.SuperUser.FirstName,
		cfg.SuperUser.LastName,
	)
	if err != nil {
		log.Printf("Note: %v", err)
	} else {
		log.Printf("Superuser created: %s (%s)", superUser.Email, superUser.Role)
	}
	
	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, cfg)
	
	// Initialize rate limiter (5 requests per minute for auth endpoints)
	authRateLimiter := middleware.NewRateLimiter(5, 1*time.Minute)
	
	// Setup routes
	mux := http.NewServeMux()
	
	// API routes
	mux.HandleFunc("/api/register", authRateLimiter.Middleware(userHandler.Register))
	mux.HandleFunc("/api/login", authRateLimiter.Middleware(userHandler.Login))
	mux.HandleFunc("/api/logout", userHandler.Logout)
	mux.HandleFunc("/api/forgot-password", authRateLimiter.Middleware(userHandler.ForgotPassword))
	mux.HandleFunc("/api/reset-password", userHandler.ResetPassword)
	mux.HandleFunc("/api/me", userHandler.GetCurrentUser)
	
	// Serve static files
	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// Serve HTML pages
	mux.HandleFunc("/", serveFile("./web/templates/login.html"))
	mux.HandleFunc("/login", serveFile("./web/templates/login.html"))
	mux.HandleFunc("/register", serveFile("./web/templates/register.html"))
	mux.HandleFunc("/forgot-password", serveFile("./web/templates/forgot-password.html"))
	mux.HandleFunc("/reset-password", serveFile("./web/templates/reset-password.html"))
	mux.HandleFunc("/dashboard", serveFile("./web/templates/dashboard.html"))
	
	// Apply global middleware
	handler := middleware.SecurityHeaders(middleware.CORS(mux))
	
	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Printf("Login at: http://localhost:%s/login", cfg.Server.Port)
	log.Printf("Superuser credentials - Email: %s, Password: %s", cfg.SuperUser.Email, cfg.SuperUser.Password)
	
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// serveFile returns a handler that serves a specific file
func serveFile(filepath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath)
	}
}
