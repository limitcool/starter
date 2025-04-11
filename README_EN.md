# Starter

[![Go Reference](https://pkg.go.dev/badge/github.com/limitcool/starter.svg)](https://pkg.go.dev/github.com/limitcool/starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/limitcool/starter)](https://goreportcard.com/report/github.com/limitcool/starter)

English | [中文](README.md)

## Features
- Provides a Gin framework project template
- Integrates GORM for ORM mapping and database operations
  - Supports PostgreSQL (using pgx driver)
  - Supports MySQL
  - Supports SQLite
  - Provides rich query option utility functions
- Integrates Viper for configuration management
- Provides common Gin middleware and tools
  - CORS middleware: Handles API cross-domain requests, implements CORS support
  - JWT parsing middleware: Parses and validates JWT Token from requests for API authentication
- Internationalization (i18n) support
  - Automatically selects language based on Accept-Language header
  - Multi-language support for error messages
  - Built-in English (en-US) and Chinese (zh-CN) translations
  - Easily extensible to support more languages
- Uses Cobra command-line framework, providing a clear subcommand structure
- Supports separation of database migration and server startup, improving startup speed
- Complete database migration system, supporting version control and rollback
- Built-in user, role, permission, and menu management system

## Quick Start

```bash
go install github.com/go-eagle/eagle/cmd/eagle@latest
eagle new <project name> -r https://github.com/limitcool/starter -b main
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

This project implements a complete database migration system for managing the creation, update, and rollback of database table structures.

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
        return tx.AutoMigrate(&model.SysUser{})
    },
    Down: func(tx *gorm.DB) error { // Down migration function
        return tx.Migrator().DropTable("sys_user")
    },
})
```

### Predefined Migrations

The system has predefined basic migration items:

1. User table (`sys_user`)
2. Role-related tables (`sys_role`, `sys_user_role`, `sys_role_menu`)
3. Permission-related tables (`sys_permission`, `sys_role_permission`)
4. Menu table (`sys_menu`)

### Adding New Migrations

To add new migrations in the `internal/migration/migrations.go` file:

1. Create a new registration function or add to an existing function
2. Ensure that version numbers follow timestamp order
3. Register in the `RegisterAllMigrations` function

```go
// Example: Adding new business table migration
func RegisterBusinessMigrations(migrator *Migrator) {
    migrator.Register(&MigrationEntry{
        Version: "202504080010",
        Name:    "create_products_table",
        Up: func(tx *gorm.DB) error {
            return tx.AutoMigrate(&model.Product{})
        },
        Down: func(tx *gorm.DB) error {
            return tx.Migrator().DropTable("products")
        },
    })
}

// Add in RegisterAllMigrations
func RegisterAllMigrations(migrator *Migrator) {
    // Existing migrations...
    RegisterBusinessMigrations(migrator)
}
```

### Migration Record Table

The system tracks the execution status of migrations through the `sys_migrations` table, containing the following fields:

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

## Logging Configuration

The project uses [charmbracelet/log](https://github.com/charmbracelet/log) as its logging library, supporting colorized console output and file output.

### Configuration Example

```yaml
Log:
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

## Permission System

The project integrates the Casbin RBAC permission system and dynamic menu system, implementing the following functions:

1. RBAC (Role-Based Access Control) permission model
   - User -> Role -> Permission
   - Supports resource-level and operation-level permission control

2. Dynamic menu system
   - Dynamically generates menus based on user roles
   - Menu items are associated with permissions
   - Supports multi-level menu tree structure

3. Permission verification middleware
   - CasbinMiddleware: Path and HTTP method-based permission control
   - PermissionMiddleware: Menu permission identifier-based permission control

4. Data table structure
   - sys_user - User table
   - sys_role - Role table
   - sys_menu - Menu table
   - sys_role_menu - Role-menu association table
   - sys_user_role - User-role association table
   - casbin_rule - Casbin rule table (automatically created)

5. API interfaces
   - Menu management: Create, update, delete, query
   - Role management: Create, update, delete, query
   - Role menu assignment
   - Role permission assignment
   - User role assignment

### Usage

1. Role and menu association:
   ```
   POST /api/v1/admin/role/menu
   {
     "role_id": 1,
     "menu_ids": [1, 2, 3]
   }
   ```

2. Role and permission association:
   ```
   POST /api/v1/admin/role/permission
   {
     "role_code": "admin",
     "object": "/api/v1/admin/user",
     "action": "GET"
   }
   ```

3. Get user menus:
   ```
   GET /api/v1/user/menus
   ```

4. Get user permissions:
   ```
   GET /api/v1/user/perms
   ```

## Database Query Options System

This project implements a complete database query options system to simplify the GORM query building process, improving code reusability and readability.

### Query Option Features

- Designed with functional option pattern
- Supports chaining multiple query conditions
- Provides a unified interface approach to handle various query scenarios
- Easy to extend and customize new query conditions

### Basic Usage

```go
// Import query options package
import "your-project/internal/pkg/options"

// Create query instance
query := options.Apply(
    DB, // *gorm.DB instance
    options.WithPage(1, 10),
    options.WithOrder("created_at", "desc"),
    options.WithLike("name", keyword),
)

// Execute query
var results []YourModel
query.Find(&results)
```

### Built-in Query Options

The system provides the following built-in query options:

#### Pagination and Sorting
- `WithPage(page, pageSize)` - Pagination query, automatically limits maximum page size
- `WithOrder(field, direction)` - Sorting query, direction supports "asc" or "desc"

#### Association Queries
- `WithPreload(relation, args...)` - Preload associations
- `WithJoin(query, args...)` - Join query
- `WithSelect(query, args...)` - Specify query fields
- `WithGroup(query)` - Group query
- `WithHaving(query, args...)` - HAVING condition query

#### Condition Filtering
- `WithWhere(query, args...)` - WHERE condition
- `WithOrWhere(query, args...)` - OR WHERE condition
- `WithLike(field, value)` - LIKE fuzzy query
- `WithExactMatch(field, value)` - Exact match query
- `WithTimeRange(field, start, end)` - Time range query
- `WithKeyword(keyword, fields...)` - Keyword search (multi-field OR condition)

#### Combined Queries
- `WithBaseQuery(tableName, status, keyword, keywordFields, createBy, startTime, endTime)` - Apply basic query conditions, combining multiple common filter conditions

### Custom Query Options

Custom query options can be easily extended:

```go
// Custom query option example
func WithCustomCondition(param string) options.Option {
    return func(db *gorm.DB) *gorm.DB {
        if param == "" {
            return db
        }
        return db.Where("custom_field = ?", param)
    }
}

// Using custom query options
query := options.Apply(
    DB,
    options.WithPage(1, 10),
    WithCustomCondition("value"),
)
```

### Using with DTOs

Query conditions can be flexibly built in combination with DTO objects:

```go
// Build query options based on BaseQuery
func BuildQueryOptions(q *request.BaseQuery, tableName string) []options.Option {
    var opts []options.Option

    // Add basic query conditions
    opts = append(opts, options.WithBaseQuery(
        tableName,
        q.Status,
        q.Keyword,
        []string{"name", "description"}, // Keyword search fields
        q.CreateBy,
        q.StartTime,
        q.EndTime,
    ))

    return opts
}

// Use in service
func (s *Service) List(query *request.YourQuery) ([]YourModel, int64, error) {
    opts := BuildQueryOptions(&query.BaseQuery, "your_table")

    // Add pagination and sorting
    opts = append(opts,
        options.WithPage(query.Page, query.PageSize),
        options.WithOrder(query.SortField, query.SortOrder),
    )

    // Apply all query options
    db := options.Apply(s.DB, opts...)

    var total int64
    db.Model(&YourModel{}).Count(&total)

    var items []YourModel
    db.Find(&items)

    return items, total, nil
}
```

## Component Access Standards

This project uses the component pattern to manage various resources, and the following access methods are recommended:

```go
// Get database connection
db := sqldb.Instance().DB()

// Get Redis client
client := redisdb.Instance().Client()

// Get MongoDB database
mongo := mongodb.Instance().DB()
```