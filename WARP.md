# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

`tableio` is a Go library for rapid prototyping of database persistence. It uses Go generics to provide a type-safe ORM-like interface that automatically generates tables from struct definitions. The library supports SQLite, MySQL, PostgreSQL, and MSSQL.

## Commands

### Testing
```bash
# Run all tests
go test -v

# Run specific test
go test -v -run TestValidateStruct

# MySQL tests require MYSQL_CONNECTION_STRING environment variable
# PostgreSQL tests require PGSQL_CONNECTION_STRING environment variable
# MSSQL tests require MSSQL_CONNECTION_STRING environment variable
```

### Development
```bash
# Install dependencies
go mod download

# Format code
go fmt ./...

# Run go vet
go vet ./...
```

## Architecture

### Core Components

**tableio.go** - Main library entry point
- `TableIO[T]` struct: Generic type wrapping database operations for type T
- `NewTableIO[T]()`: Constructor that validates struct, establishes DB connection, caches field metadata
- CRUD operations: `Insert()`, `InsertMany()`, `All()`
- Schema operations: `CreateTableIfNotExists()`, `DeleteTableIfExists()`

**reflectx/** - Reflection utilities package
- `reflect.go`: Core reflection operations for extracting struct metadata, generating SQL DDL, and converting struct values to SQL-safe strings
- `strings.go`: String manipulation utilities (snake_case conversion, suffix trimming)
- `FieldInfo` struct: Metadata container for struct field name, type, and db tag

### Key Design Patterns

**Struct Tag System**
TableIO uses `tableio:"..."` struct tags to define field behavior:
- `pk` or `primarykey` - Marks field as primary key
- `auto` or `autoincrement` - Database auto-generates this field (IDENTITY/SERIAL/AUTO_INCREMENT)
- `unique` - Adds UNIQUE constraint
- `required` or `notnull` - Field is NOT NULL

Example:
```go
type Person struct {
    ID      int64  `tableio:"pk,auto"`
    Name    string `tableio:"unique,required"`
    Age     int
    Address Address
}
```

Fields without tags default to nullable with no constraints.

**Table Naming Convention**
Table names are automatically derived from struct type names:
1. Extract type name from struct
2. Pluralize (using `github.com/gertd/go-pluralize`)
3. Convert to snake_case (using `github.com/iancoleman/strcase`)

Example: `Person` struct → `people` table, `AzureCloudspace` struct → `azure_cloudspaces` table

**Field Mapping**
- All struct fields are mapped to DB columns
- Field names become column names (preserving case)
- Type mapping: `string` → VARCHAR(255), `int/int32/int64` → INTEGER, other types → JSON/TEXT/NVARCHAR(MAX) (database-specific)
- Fields marked with `auto` get database-specific auto-increment: AUTO_INCREMENT (MySQL), SERIAL (PostgreSQL), AUTOINCREMENT (SQLite), IDENTITY(1,1) (MSSQL)

**Nested Struct Handling**
Nested structs (like `Address` in `Person`) are serialized to JSON and stored in JSON/TEXT columns.

**SQL Generation**
- SELECT statements use cached `selectList` (all fields)
- INSERT statements use cached `insertList` (excludes auto-increment fields)
- Field lists and SQL DDL are generated once during `NewTableIO()` construction for performance
- CREATE/DROP TABLE statements are database-specific (MSSQL uses `IF NOT EXISTS` with `BEGIN/END` blocks)

### Dependencies

External:
- `github.com/denisenkom/go-mssqldb` - MSSQL driver
- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/gertd/go-pluralize` - Table name pluralization
- `github.com/iancoleman/strcase` - Snake case conversion
- `github.com/katasec/utils/errx` - Error handling utilities

## Important Notes

- The `All()` method currently prints JSON to stdout - consider this for debugging vs production use
- Error handling uses `errx.PanicOnError()` - operations panic rather than returning errors in most cases
- Connection strings are database-specific - refer to README.md examples for format
- Tests require environment variables for MySQL and PostgreSQL connections, MSSQL test has hardcoded connection
- No struct validation is performed - use struct tags to define constraints as needed
