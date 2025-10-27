# Angidi - User Authentication Microservice

A simple yet secure user authentication microservice built with Go, implementing basic authentication with support for future extensibility to advanced authentication mechanisms.

## Features

### ✅ Implemented Features

- **User Registration**: Self-service account creation with automatic "Customer/Viewer" role assignment
- **User Login**: Email and password-based authentication
- **Password Security**: 
  - Bcrypt password hashing
  - Strong password validation (minimum 8 characters with uppercase, lowercase, and numbers)
  - Account lockout after 5 failed login attempts (30-minute lockout)
- **Forgot Password**: Password reset functionality with time-limited tokens (1-hour expiry)
- **User Roles**: SuperUser, Admin, Viewer, and Customer roles
- **User Status**: Active, Pending Verification, and Disabled statuses
- **Remember Me**: Persistent sessions with configurable cookie lifetime
- **Rate Limiting**: Protection against brute-force attacks (5 requests per minute on auth endpoints)
- **Security Headers**: CORS, XSS protection, and other security headers
- **Superuser Bootstrap**: Automatic superuser creation on startup
- **In-Memory Storage**: Thread-safe in-memory data store (ready for PostgreSQL migration)

### 🎨 Frontend Features

- Login page with "Remember Me" option
- Registration page with password confirmation
- Forgot password page
- Password reset page
- User dashboard displaying user information
- Responsive design with modern gradient UI
- Client-side validation

### 🏗️ Architecture Highlights

- **Strategy Pattern**: Pluggable authentication mechanisms (currently supporting Basic Auth, ready for JWT, 2FA, OAuth)
- **Clean Architecture**: Separation of concerns with domain, repository, service, and handler layers
- **Middleware Support**: Rate limiting, CORS, and security headers
- **Thread-Safe**: Concurrent-safe in-memory repository with mutex locks

## Getting Started

### Prerequisites

- Go 1.24 or higher
- No external dependencies for basic operation

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yesoreyeram/angidi.git
cd angidi
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o bin/user-service ./cmd/user-service
```

### Running the Service

```bash
./bin/user-service
```

The service will start on `http://localhost:8080` by default.

**Default Superuser Credentials:**
- Email: `admin@angidi.com`
- Password: `Admin@123`

### Configuration

The service can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Server port |
| `SERVER_HOST` | `0.0.0.0` | Server host |
| `SUPERUSER_EMAIL` | `admin@angidi.com` | Superuser email |
| `SUPERUSER_PASSWORD` | `Admin@123` | Superuser password |
| `SUPERUSER_FIRST_NAME` | `Super` | Superuser first name |
| `SUPERUSER_LAST_NAME` | `Admin` | Superuser last name |
| `SESSION_SECRET` | `change-me-in-production` | Session secret key |
| `SESSION_COOKIE_MAX_AGE` | `86400` | Session cookie max age (seconds) |

## API Endpoints

### Authentication Endpoints

#### Register a New User
```bash
POST /api/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass@123",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Login
```bash
POST /api/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass@123",
  "remember_me": true
}
```

#### Logout
```bash
POST /api/logout
```

#### Forgot Password
```bash
POST /api/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

#### Reset Password
```bash
POST /api/reset-password
Content-Type: application/json

{
  "token": "reset_token_here",
  "new_password": "NewSecurePass@123"
}
```

#### Get Current User
```bash
GET /api/me
```

## Project Structure

```
angidi/
├── cmd/
│   └── user-service/        # Main application entry point
│       └── main.go
├── internal/
│   ├── config/              # Configuration management
│   ├── domain/              # Domain models and DTOs
│   ├── errors/              # Application errors
│   ├── handler/             # HTTP handlers
│   ├── middleware/          # HTTP middleware (CORS, rate limiting)
│   ├── repository/          # Data access layer
│   └── service/             # Business logic layer
├── pkg/
│   └── auth/                # Authentication strategies and utilities
├── web/
│   ├── static/              # Static assets (CSS, JS)
│   │   ├── css/
│   │   └── js/
│   └── templates/           # HTML templates
├── docs/                    # Project documentation
└── go.mod
```

## Security Features

- **Password Hashing**: Using bcrypt with default cost factor
- **Account Lockout**: Automatic lockout after 5 failed login attempts
- **Rate Limiting**: 5 requests per minute on authentication endpoints
- **Secure Sessions**: HttpOnly and SameSite cookies
- **Security Headers**: X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, HSTS
- **Password Strength**: Enforced minimum complexity requirements
- **Token Expiry**: Time-limited password reset tokens (1 hour)

## Future Enhancements

The architecture is designed to support the following future features:

### Authentication Mechanisms
- **JWT Authentication**: Token-based authentication
- **2FA/MFA**: Two-factor authentication using TOTP
- **OAuth 2.0**: Social login (Google, GitHub, etc.)
- **Biometric Auth**: Fingerprint, Face ID support

### Data Persistence
- **PostgreSQL**: Relational database for user data
- **Redis**: Session management and caching
- **MongoDB**: Flexible schema for user preferences

### Additional Features
- Email verification for new accounts
- Security questions for account recovery
- Session management (view/revoke active sessions)
- Audit logging for security events
- User profile management
- Role and permission management UI

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/user-service ./cmd/user-service
```

### Adding a New Authentication Strategy

1. Create a new strategy file in `pkg/auth/`:
```go
type MyAuthStrategy struct {
    // ... dependencies
}

func (s *MyAuthStrategy) Name() string {
    return "myauth"
}

func (s *MyAuthStrategy) Authenticate(credentials interface{}) (*domain.User, error) {
    // ... implementation
}
```

2. Register the strategy in `cmd/user-service/main.go`:
```go
myAuth := auth.NewMyAuthStrategy(userRepo)
authMgr.RegisterStrategy(myAuth)
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue on the GitHub repository
