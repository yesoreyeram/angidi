# PostgreSQL Migration Guide

This document outlines how to migrate from the in-memory data store to PostgreSQL persistence.

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('superuser', 'admin', 'viewer', 'customer')),
    status VARCHAR(50) NOT NULL CHECK (status IN ('active', 'pending_verification', 'disabled')),
    failed_login_count INTEGER DEFAULT 0,
    last_failed_login TIMESTAMP,
    password_reset_token VARCHAR(255),
    password_reset_expiry TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_password_reset_token ON users(password_reset_token) WHERE password_reset_token IS NOT NULL;
```

## Implementation Steps

### 1. Add PostgreSQL Dependencies

Update `go.mod`:

```bash
go get github.com/lib/pq
go get github.com/jmoiron/sqlx
```

### 2. Create PostgreSQL Repository

Create `internal/repository/postgres_user_repository.go`:

```go
package repository

import (
    "database/sql"
    "time"
    
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    
    "github.com/yesoreyeram/angidi/internal/domain"
    "github.com/yesoreyeram/angidi/internal/errors"
)

type PostgresUserRepository struct {
    db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
    query := `
        INSERT INTO users (id, email, password_hash, first_name, last_name, role, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
    _, err := r.db.Exec(query,
        user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
        user.Role, user.Status, user.CreatedAt, user.UpdatedAt,
    )
    
    if err != nil {
        // Check for unique constraint violation
        if strings.Contains(err.Error(), "duplicate key") {
            return errors.ErrUserAlreadyExists
        }
        return err
    }
    
    return nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*domain.User, error) {
    var user domain.User
    query := `SELECT * FROM users WHERE email = $1`
    
    err := r.db.Get(&user, query, email)
    if err == sql.ErrNoRows {
        return nil, errors.ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

// Implement other methods...
```

### 3. Add Database Configuration

Update `internal/config/config.go`:

```go
type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func Load() *Config {
    return &Config{
        // ... existing config
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 5432),
            User:     getEnv("DB_USER", "angidi"),
            Password: getEnv("DB_PASSWORD", ""),
            DBName:   getEnv("DB_NAME", "angidi"),
            SSLMode:  getEnv("DB_SSL_MODE", "disable"),
        },
    }
}
```

### 4. Initialize Database Connection

Update `cmd/user-service/main.go`:

```go
import (
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func initDB(cfg *config.Config) (*sqlx.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.DBName,
        cfg.Database.SSLMode,
    )
    
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Set connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    return db, nil
}

func main() {
    cfg := config.Load()
    
    // Initialize database
    db, err := initDB(cfg)
    if err != nil {
        log.Fatalf("Database initialization failed: %v", err)
    }
    defer db.Close()
    
    // Use PostgreSQL repository instead of in-memory
    userRepo := repository.NewPostgresUserRepository(db)
    
    // ... rest of the code
}
```

### 5. Database Migrations

Create a migration tool or use a library like `golang-migrate/migrate`:

```bash
go get -u github.com/golang-migrate/migrate/v4
```

Create migration files:
- `migrations/000001_create_users_table.up.sql`
- `migrations/000001_create_users_table.down.sql`

### 6. Environment Variables

Add to `.env` or export:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=angidi
DB_PASSWORD=your_secure_password
DB_NAME=angidi
DB_SSL_MODE=require  # Use 'require' in production
```

### 7. Docker Setup (Optional)

Create `docker-compose.yml` for local development:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: angidi
      POSTGRES_PASSWORD: angidi_dev_password
      POSTGRES_DB: angidi
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## Testing

### Unit Tests with Mock Database

```go
func TestUserRepository(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to create mock: %v", err)
    }
    defer db.Close()
    
    sqlxDB := sqlx.NewDb(db, "sqlmock")
    repo := NewPostgresUserRepository(sqlxDB)
    
    // Test cases...
}
```

### Integration Tests

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Test with real database...
}
```

## Performance Considerations

1. **Indexes**: Add indexes on frequently queried columns (email, status, created_at)
2. **Connection Pooling**: Configure appropriate pool sizes based on load
3. **Prepared Statements**: Use for frequently executed queries
4. **Read Replicas**: Consider for read-heavy workloads
5. **Caching**: Add Redis for session management and frequently accessed data

## Monitoring

Add monitoring for:
- Connection pool metrics
- Query performance
- Database size and growth
- Slow query log
- Replication lag (if using replicas)

## Backup Strategy

1. Automated daily backups
2. Point-in-time recovery capability
3. Regular restore testing
4. Off-site backup storage

## Security

1. Use SSL/TLS for database connections in production
2. Encrypt sensitive fields (consider using pgcrypto)
3. Implement row-level security policies
4. Regular security audits
5. Principle of least privilege for database users
