# Starter (Lite Version)

[![Go Reference](https://pkg.go.dev/badge/github.com/limitcool/starter.svg)](https://pkg.go.dev/github.com/limitcool/starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/limitcool/starter)](https://goreportcard.com/report/github.com/limitcool/starter)

English | [中文](README.md)

> This is the lightweight version of the Starter framework, designed with a single-user mode, suitable for rapid development and simple application scenarios. If you need more enterprise-level features, please check the `enterprise` branch.

## Features
- Provides a lightweight Gin framework project template
- Uses simple manual dependency injection, creating clearer code structure
- Adopts simplified architecture design, focused on rapid development
- Integrates GORM for ORM mapping and database operations
  - Supports PostgreSQL (using pgx driver)
  - Supports MySQL
  - Supports SQLite
- Integrates Viper for configuration management
- Provides common Gin middleware and tools
  - CORS middleware: Handles API cross-domain requests, implements CORS support
  - JWT parsing middleware: Parses and validates JWT Token from requests for API authentication
- Internationalization (i18n) support
  - Automatically selects language based on Accept-Language header
  - Multi-language support for error messages
- Uses Cobra command-line framework, providing a clear subcommand structure
- Supports separation of database migration and server startup, improving startup speed
- Simplified user management system, using a single user table with IsAdmin field to distinguish administrators
- Optimized error handling system, supporting error codes and multilingual error messages

## Architecture Design

The project adopts a simplified architecture design, combined with the Uber fx dependency injection framework, implementing a clear code structure:

### 1. Simplified Layered Architecture

- **Model Layer**: Defines data models and database table structures
- **Handler Layer**: Handles HTTP requests and responses, directly interacts with the database
- **Router Layer**: Defines API routes, depends on the Handler layer

### 2. Dependency Injection

The project uses the Uber fx framework to implement dependency injection, with dependencies injected through constructors:

```go
// Model layer
func NewUserRepo(db *gorm.DB) *UserRepo {
    // ...
}

// Handler layer
func NewUserHandler(db *gorm.DB, config *configs.Config) *UserHandler {
    // ...
}

// Router layer
func NewRouter(userHandler *handler.UserHandler) *gin.Engine {
    // ...
}
```

### 3. Application Container Management

Uses application container to manage component lifecycles, ensuring proper initialization and cleanup:

```go
type App struct {
    config   *configs.Config
    db       *gorm.DB
    handlers *Handlers
    router   *gin.Engine
    server   *http.Server
}

func New(config *configs.Config) (*App, error) {
    app := &App{config: config}

    // Initialize components in order
    if err := app.initDatabase(); err != nil {
        return nil, err
    }
    // ... other component initialization

    return app, nil
}
```

## Quick Start

```bash
go install github.com/go-eagle/eagle/cmd/eagle@latest
eagle new <project name> -r https://github.com/limitcool/starter -b lite
```

## Usage

The application uses the Cobra command-line framework, providing a clearer subcommand structure.

### Basic Commands

```bash
# View help information
./<app-name> --help

# View version information
./<app-name> version

# Start the server
./<app-name> server

# Execute database migration
./<app-name> migrate
```

### Server Command

The server command is used to start the HTTP service:

```bash
# Start the server with default configuration
./<app-name> server

# Start the server with a specified port
./<app-name> server --port 9000

# Start the server with a specified configuration file
./<app-name> server --config custom.yaml
```

### Database Migration Commands

Database migration commands are used to initialize or update the database structure:

```bash
# Execute database migration
./<app-name> migrate

# Execute migration with a specified configuration file
./<app-name> migrate --config prod.yaml

# Clear the database before migration (dangerous operation)
./<app-name> migrate --fresh

# Rollback the last batch of database migrations
./<app-name> migrate rollback

# Display database migration status
./<app-name> migrate status

# Reset all database migrations
./<app-name> migrate reset
```

## Database Migration System

The Lite version implements a concise database migration system for managing the creation and update of database table structures.

### Migration System Features

- Supports executing migrations in version number order
- Tracks executed migration records
- Supports transactional migrations, ensuring data consistency
- Provides up and down migration functions
- Supports batch rollback and complete reset

### Migration File Structure

Migrations are defined in the `internal/migration/migrations.go` file, following this structure:

```go
migrator.Register(&MigrationEntry{
    Version: "202504080001",        // Version number format: YearMonthDaySerialNumber
    Name:    "create_users_table",  // Migration name
    Up: func(tx *gorm.DB) error {   // Up migration function
        return tx.AutoMigrate(&model.User{})
    },
    Down: func(tx *gorm.DB) error { // Down migration function
        return tx.Migrator().DropTable("users")
    },
})
```

### Predefined Migrations

The Lite version has predefined basic migration items:

1. User table (`users`)
2. File table (`files`)

### Adding New Migrations

To add new migrations, in the `internal/migration/migrations.go` file:

1. Create a new registration function or add to an existing function
2. Ensure that version numbers follow timestamp order
3. Use the `RegisterMigration` function to register

```go
// Example: Adding new business table migration
RegisterMigration("create_products_table",
    // Up migration function
    func(tx *gorm.DB) error {
        return tx.AutoMigrate(&model.Product{})
    },
    // Down migration function
    func(tx *gorm.DB) error {
        return tx.Migrator().DropTable("products")
    },
)
```

### Migration Record Table

The system tracks the execution status of migrations through the `migrations` table, containing the following fields:

- `id`: Auto-increment primary key
- `version`: Migration version number (unique index)
- `name`: Migration name
- `created_at`: Execution time
- `batch`: Batch number (for rollback)

## Environment Configuration

Specify the running environment via the `APP_ENV` environment variable, or directly specify the configuration file via the `--config` flag:

- `APP_ENV=dev` or `APP_ENV=development` - Development environment (default)
- `APP_ENV=test` or `APP_ENV=testing` - Testing environment
- `APP_ENV=prod` or `APP_ENV=production` - Production environment

Examples:
```bash
# Run the server in development environment
APP_ENV=dev ./<app-name> server

# Execute database migration in production environment
APP_ENV=prod ./<app-name> migrate
```

## Configuration File

Configuration files are automatically loaded according to the running environment:

- `dev.yaml` - Development environment configuration
- `test.yaml` - Testing environment configuration
- `prod.yaml` - Production environment configuration
- `example.yaml` - Example configuration (for version control)

Configuration files can be placed in the following locations (in order of lookup):
1. Current working directory (project root)
2. `configs/` directory

When using for the first time, please copy the example configuration and rename it according to the environment:

```bash
# Development environment (in root directory)
cp example.yaml ./dev.yaml

# Or in configs directory
cp example.yaml configs/dev.yaml

# Production environment
cp example.yaml configs/prod.yaml
```

The application will automatically find and load the corresponding configuration file based on the `APP_ENV` environment variable. For example, when `APP_ENV=dev`, it will look for configuration files in the following order:
1. `./dev.yaml` (current directory)
2. `./configs/dev.yaml` (configs directory)

If the corresponding configuration file cannot be found, the application will not start.

## Internationalization (i18n) Support

The system has built-in internationalization support, which can automatically switch languages based on client requests.

### Configuring Internationalization

Set internationalization options in the configuration file:

```yaml
I18n:
  Enabled: true                # Whether to enable internationalization
  DefaultLanguage: en-US       # Default language
  SupportLanguages:            # List of supported languages
    - zh-CN
    - en-US
  ResourcesPath: locales       # Language resource file path
```

### Language Resource Files

Language resource files are located in the `locales` directory, in JSON format:

- `locales/en-US.json` - English resources
- `locales/zh-CN.json` - Chinese resources

Example language file content:

```json
{
  "error.success": "Success",
  "error.common.invalid_params": "Invalid request parameters",
  "error.user.user_not_found": "User not found"
}
```

### Usage

1. **Automatic translation of API responses**:
   - The system automatically selects the language based on the `Accept-Language` request header
   - API error responses will return translated text according to the set language

2. **Client request examples**:
   ```bash
   # Request English response
   curl -X POST "http://localhost:8080/api/v1/user/login" \
        -H "Accept-Language: en-US" \
        -H "Content-Type: application/json" \
        -d '{"username": "test", "password": "wrong"}'

   # Request Chinese response
   curl -X POST "http://localhost:8080/api/v1/user/login" \
        -H "Accept-Language: zh-CN" \
        -H "Content-Type: application/json" \
        -d '{"username": "test", "password": "wrong"}'
   ```

3. **Adding new error code translations**:
   - Define errors in `tools/errorgen/error_codes.md`
   - Run the error code generator: `go run tools/errorgen/main.go tools/errorgen/error_codes.md internal/pkg/errorx/code_gen.go`
   - Add corresponding translations in language files (`locales/en-US.json` and `locales/zh-CN.json`)

4. **Adding support for a new language**:
   - Create a new language file, e.g., `locales/fr-FR.json`
   - Add the language to the `SupportLanguages` list in the configuration
   - Restart the application to make the configuration effective

## Error Handling System

The project implements a complete error handling system, including error codes, error wrapping, and multilingual error messages.

### Error Handling Features

- Unified error code definition and management
- Error wrapping, preserving complete error chains and stack information
- Multilingual error message support
- Distinction between internal errors and user-visible errors

### Error Handling Best Practices

- Repository layer: Returns specific errors, does not log
- Service layer: Wraps errors, adds business context, does not log
- Controller layer: Converts to user-friendly error responses, logs complete error information

### Usage Example

```go
// Repository layer
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
    var user model.User
    if err := r.DB.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errorx.NewError(errorx.ErrorUserNotFoundCode, "User not found")
        }
        return nil, err
    }
    return &user, nil
}

// Service layer
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, errorx.WrapError(err, fmt.Sprintf("Failed to get user with ID %d", id))
    }
    return user, nil
}

// Controller layer
func (c *UserController) GetUser(ctx *gin.Context) {
    id := cast.ToInt64(ctx.Param("id"))
    user, err := c.userService.GetUserByID(ctx, id)
    if err != nil {
        logger.Error("Failed to get user", "error", err, "id", id)
        response.Error(ctx, err)
        return
    }
    response.Success(ctx, user)
}
```

## Logging Configuration

The project supports multiple logging libraries, including [charmbracelet/log](https://github.com/charmbracelet/log) and [uber-go/zap](https://github.com/uber-go/zap), which can be switched via configuration.

### Configuration Example

```yaml
Log:
  Driver: charm               # Log driver: charm, zap
  Level: info                 # Log level: debug, info, warn, error
  Output: [console, file]     # Output methods: console, file
  Format: text                # Log format: text, json
  FileConfig:
    Path: ./logs/app.log      # Log file path
    MaxSize: 100              # Maximum size of each log file (MB)
    MaxAge: 7                 # Days to retain log files
    MaxBackups: 10            # Maximum number of old log files to retain
    Compress: true            # Whether to compress old log files
```

### Log Drivers

- `charm`: Uses the charmbracelet/log library, supports colorized console output, suitable for development environments
- `zap`: Uses the uber-go/zap library, high-performance structured logging, suitable for production environments

### Log Levels

- `debug`: Debug information, including detailed development debugging information
- `info`: General information, default level
- `warn`: Warning information, information that needs attention
- `error`: Error information, errors that affect normal program operation

### Log Formats

- `text`: Plain text format, suitable for human reading (default)
- `json`: JSON structured format, suitable for machine parsing and log system collection

### Output Methods

- `console`: Output to console, supporting colorized output
- `file`: Output to file, supporting automatic splitting by size, automatic cleaning, and compression

Multiple output methods can be configured simultaneously, and logs will be output to all configured targets. If output is not configured, it defaults to console only.

### File Output Configuration

- `Path`: Log file path
- `MaxSize`: Maximum size of a single log file (MB), automatically split after exceeding
- `MaxAge`: Number of days to retain log files, automatically deleted after exceeding
- `MaxBackups`: Number of old log files to retain
- `Compress`: Whether to compress old log files

### Usage Example

The logging library provides a unified interface, so the usage is consistent regardless of which driver is used:

```go
// Import the logger package
import "github.com/limitcool/starter/internal/pkg/logger"

// Log at different levels
func example() {
    // Log info message
    logger.Info("This is an info message", "user", "admin", "action", "login")

    // Log warning message
    logger.Warn("This is a warning message", "memory", "90%")

    // Log error message
    err := someFunction()
    if err != nil {
        logger.Error("Operation failed", "error", err, "operation", "someFunction")
    }

    // Log debug message
    logger.Debug("Detailed debug information", "request", req, "response", resp)
}
```

To switch the log driver, you only need to modify `Log.Driver` in the configuration file, with no need to change your code.

## gRPC Support

The project integrates gRPC support, which can run in parallel with HTTP services, providing high-performance RPC services.

### gRPC Features

- Can enable/disable gRPC service via configuration
- Supports gRPC health check and reflection services
- Shares business logic with HTTP services
- Uses Protocol Buffers to define APIs

### Configuration Example

```yaml
GRPC:
  Enabled: true                # Whether to enable gRPC service
  Port: 9000                  # gRPC service port
  HealthCheck: true           # Whether to enable health check service
  Reflection: true            # Whether to enable reflection service
```

### Usage

1. **Define Proto Files**

   Proto files are defined in the `internal/proto/v1` directory, and the generated code is in the `internal/proto/gen/v1` directory.

   ```protobuf
   // internal/proto/v1/system.proto
   syntax = "proto3";

   package internal.proto.v1;

   option go_package = "internal/proto/gen/v1;protov1";

   // SystemService system service
   service SystemService {
     // GetSystemInfo get system information
     rpc GetSystemInfo(SystemInfoRequest) returns (SystemInfoResponse) {}
   }

   // SystemInfoRequest system information request
   message SystemInfoRequest {
     // Request ID
     string request_id = 1;
   }

   // SystemInfoResponse system information response
   message SystemInfoResponse {
     // Application name
     string app_name = 1;
     // Application version
     string version = 2;
     // Running mode
     string mode = 3;
     // Server time
     int64 server_time = 4;
   }
   ```

2. **Generate gRPC Code**

   Use the proto command in Makefile to generate gRPC code:

   ```bash
   make proto
   ```

3. **Implement gRPC Controllers**

   Create gRPC controllers in the `internal/controller` directory, using the `_grpc` suffix to distinguish:

   ```go
   // internal/controller/system_grpc.go
   package controller

   import (
       "context"
       "time"

       pb "github.com/limitcool/starter/internal/proto/gen/v1"
       // ...
   )

   // SystemGRPCController gRPC system controller
   type SystemGRPCController struct {
       pb.UnimplementedSystemServiceServer
       // ...
   }

   // GetSystemInfo get system information
   func (c *SystemGRPCController) GetSystemInfo(ctx context.Context, req *pb.SystemInfoRequest) (*pb.SystemInfoResponse, error) {
       // Implement business logic
       return &pb.SystemInfoResponse{
           AppName:    c.config.App.Name,
           Version:    "1.0.0",
           Mode:       c.config.App.Mode,
           ServerTime: time.Now().Unix(),
       }, nil
   }
   ```

4. **Register gRPC Services**

   Register gRPC controllers in `internal/controller/module.go`:

   ```go
   // Initialize gRPC server in application container
   func (a *App) initGRPCServer() error {
       // Create gRPC server
       grpcServer := grpc.NewServer()

       // Register services
       systemController := NewSystemGRPCController(a.config)
       pb.RegisterSystemServiceServer(grpcServer, systemController)

       a.grpcServer = grpcServer
       return nil
   }
   ```

5. **Use gRPC Client**

   ```go
   // Create gRPC connection
   conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
   if err != nil {
       log.Fatalf("Failed to connect: %v", err)
   }
   defer conn.Close()

   // Create client
   client := pb.NewSystemServiceClient(conn)

   // Call service
   resp, err := client.GetSystemInfo(context.Background(), &pb.SystemInfoRequest{
       RequestId: "test-request",
   })
   ```

## Permission System

The Lite version adopts a simplified permission system, based on the user's `is_admin` field for permission control:

### Permission Control Middleware

The project provides three types of permission control middleware:

1. **AdminCheck**: Checks if the user is an administrator
   - Based on the `is_admin` field in JWT for quick checking
   - Suitable for admin-exclusive interfaces

2. **UserCheck**: Checks if the user is logged in
   - Only verifies that the user is logged in, does not check user type
   - Suitable for interfaces that require login but do not restrict user type

3. **RegularUserCheck**: Checks if the user is a regular user
   - Ensures the user is not an administrator
   - Suitable for interfaces that only allow regular users to access

### Usage

Using middleware in route definitions:

```go
// Admin interfaces
adminGroup := router.Group("/api/v1/admin")
adminGroup.Use(middleware.AdminCheck())
{
    adminGroup.GET("/users", handler.ListUsers)
    // Other admin interfaces...
}

// Regular user interfaces
userGroup := router.Group("/api/v1/user")
userGroup.Use(middleware.UserCheck())
{
    userGroup.GET("/profile", handler.GetUserProfile)
    // Other user interfaces...
}

// Regular-user-only interfaces
regularUserGroup := router.Group("/api/v1/regular")
regularUserGroup.Use(middleware.RegularUserCheck())
{
    regularUserGroup.POST("/feedback", handler.SubmitFeedback)
    // Other regular-user-only interfaces...
}
```

## Database Operations

The Lite version adopts a simplified database operation approach, providing database operation methods directly in the Model layer:

### Model Layer Design

In the Lite version, the Model layer directly provides database operation methods, simplifying the code structure:

```go
// User model
type User struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    Username  string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
    Password  string    `gorm:"size:100;not null" json:"-"`
    Nickname  string    `gorm:"size:50" json:"nickname"`
    Email     string    `gorm:"size:100" json:"email"`
    Avatar    string    `gorm:"size:255" json:"avatar"`
    IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// UserRepo database operations for users
type UserRepo struct {
    DB *gorm.DB
}

// NewUserRepo creates a user repository
func NewUserRepo(db *gorm.DB) *UserRepo {
    return &UserRepo{DB: db}
}

// GetByID gets a user by ID
func (r *UserRepo) GetByID(ctx context.Context, id uint) (*User, error) {
    var user User
    if err := r.DB.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errorx.ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

// GetByUsername gets a user by username
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
    var user User
    if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errorx.ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

// Create creates a user
func (r *UserRepo) Create(ctx context.Context, user *User) error {
    return r.DB.Create(user).Error
}

// Update updates a user
func (r *UserRepo) Update(ctx context.Context, user *User) error {
    return r.DB.Save(user).Error
}

// Delete deletes a user
func (r *UserRepo) Delete(ctx context.Context, id uint) error {
    return r.DB.Delete(&User{}, id).Error
}

// List gets a list of users
func (r *UserRepo) List(ctx context.Context, page, pageSize int) ([]User, int64, error) {
    var users []User
    var total int64

    r.DB.Model(&User{}).Count(&total)

    offset := (page - 1) * pageSize
    if err := r.DB.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, total, nil
}
```

### Using in Handler Layer

In the Handler layer, directly use the methods provided by the Model layer:

```go
// UserHandler user handler
type UserHandler struct {
    userRepo *model.UserRepo
    config   *configs.Config
}

// NewUserHandler creates a user handler
func NewUserHandler(userRepo *model.UserRepo, config *configs.Config) *UserHandler {
    return &UserHandler{
        userRepo: userRepo,
        config:   config,
    }
}

// GetUser gets user information
func (h *UserHandler) GetUser(c *gin.Context) {
    id := cast.ToUint(c.Param("id"))

    user, err := h.userRepo.GetByID(c.Request.Context(), id)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, user)
}
```

## Error Handling

The Lite version provides a concise yet powerful error handling system, supporting error codes and multilingual error messages.

### Error Handling Features

- Unified error code definition and management
- Error wrapping, preserving complete error chains
- Multilingual error message support
- Distinction between internal errors and user-visible errors

### Error Definition

Errors are defined in the `internal/pkg/errorx` package:

```go
// Error code definitions
const (
    // Common error codes
    ErrorSuccess       = 0    // Success
    ErrorUnknown       = 1000 // Unknown error
    ErrorInvalidParams = 1001 // Invalid parameters
    ErrorNotFound      = 1002 // Resource not found
    ErrorDatabase      = 1003 // Database error

    // User-related error codes
    ErrorUserNotFound     = 2000 // User not found
    ErrorUserAlreadyExist = 2001 // User already exists
    ErrorUserAuthFailed   = 2002 // User authentication failed
    ErrorUserNoLogin      = 2003 // User not logged in
    ErrorAccessDenied     = 2004 // Access denied
)

// Error custom error type
type Error struct {
    Code    int    // Error code
    Message string // Error message
    Err     error  // Original error
}

// NewError creates a new error
func NewError(code int, msg string) *Error {
    return &Error{
        Code:    code,
        Message: msg,
    }
}

// WithMsg sets the error message
func (e *Error) WithMsg(msg string) *Error {
    return &Error{
        Code:    e.Code,
        Message: msg,
        Err:     e.Err,
    }
}

// WithError wraps the original error
func (e *Error) WithError(err error) *Error {
    return &Error{
        Code:    e.Code,
        Message: e.Message,
        Err:     err,
    }
}
```

### Usage Example

```go
// In Model layer
func (r *UserRepo) GetByID(ctx context.Context, id uint) (*User, error) {
    var user User
    if err := r.DB.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errorx.NewError(errorx.ErrorUserNotFound, "User not found")
        }
        return nil, errorx.NewError(errorx.ErrorDatabase, "Database error").WithError(err)
    }
    return &user, nil
}

// In Handler layer
func (h *UserHandler) GetUser(c *gin.Context) {
    id := cast.ToUint(c.Param("id"))

    user, err := h.userRepo.GetByID(c.Request.Context(), id)
    if err != nil {
        logger.Error("Failed to get user", "error", err, "id", id)
        response.Error(c, err)
        return
    }

    response.Success(c, user)
}
```

### Unified Response Format

All API responses use a unified format:

```go
// Success response
{
    "code": 0,
    "message": "success",
    "data": {
        // Response data
    }
}

// Error response
{
    "code": 1001,
    "message": "Invalid parameters",
    "data": null
}
```
