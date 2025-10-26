# bunhelpers

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A collection of powerful helper functions and utilities for [Bun ORM](https://github.com/uptrace/bun), designed to streamline database operations and reduce boilerplate code in Go applications.

## Features

- 🔍 **Query Selectors**: Composable query builders with type-safe WHERE conditions
- 🔄 **Transaction Management**: Context-aware transaction handling with nested transaction support
- 📊 **Filtering Utilities**: Pre-built `Where` struct for common filtering patterns
- 🏗️ **Querier Interface**: Context-aware query builders that automatically use transactions from context
- ❌ **Error Handling**: Specialized error checkers for common database errors

## Installation

```bash
go get github.com/uagolang/bunhelpers
```

## Quick Start

```go
package main

import (
    "context"
    "database/sql"
    
    "github.com/uagolang/bunhelpers"
    "github.com/uptrace/bun"
    "github.com/uptrace/bun/dialect/pgdialect"
    "github.com/uptrace/bun/driver/pgdriver"
)

func main() {
    // Setup Bun
    sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://...")))
    db := bun.NewDB(sqldb, pgdialect.New())
    defer db.Close()
    
    ctx := context.Background()
    
    // Use selectors
    var users []User
    err := db.NewSelect().
        Model(&users).
        Apply(bunhelpers.WhereEqual("status", "active")).
        Apply(bunhelpers.WhereContains("email", "example.com")).
        Scan(ctx)
}
```

## Documentation

### 1. Query Selectors

Build complex queries using composable selector functions:

#### Basic Selectors

```go
// Equality checks
query.Apply(bunhelpers.WhereEqual("status", "active"))
query.Apply(bunhelpers.WhereNotEqual("role", "guest"))

// NULL checks
query.Apply(bunhelpers.WhereNull("deleted_at"))
query.Apply(bunhelpers.WhereNotNull("verified_at"))

// IN clauses
query.Apply(bunhelpers.WhereIn("id", []string{"1", "2", "3"}))
query.Apply(bunhelpers.WhereNotIn("status", []string{"banned", "suspended"}))

// String matching (case-insensitive)
query.Apply(bunhelpers.WhereContains("name", "john"))  // ILIKE '%john%'
query.Apply(bunhelpers.WhereBegins("email", "admin")) // ILIKE 'admin%'
query.Apply(bunhelpers.WhereEnds("domain", ".com"))   // ILIKE '%.com'

// Time-based queries
query.Apply(bunhelpers.WhereBefore("created_at", time.Now()))
query.Apply(bunhelpers.WhereAfter("updated_at", lastWeek))
```

#### JSONB Selectors

Work with PostgreSQL JSONB columns safely:

```go
// Simple JSONB field equality
query.Apply(bunhelpers.WhereJsonbEqual("metadata", "status", "verified"))
// Generates: metadata->>'status' = 'verified'

// Nested JSONB path equality
query.Apply(bunhelpers.WhereJsonbPathEqual("data", []string{"user", "profile", "name"}, "John"))
// Generates: data->'user'->'profile'->>'name' = 'John'

// JSONB array of objects - find object with matching field
query.Apply(bunhelpers.WhereJsonbObjectsArrayKeyValueEqual("tags", "items", "id", "123"))

// Nested array search
query.Apply(bunhelpers.WhereJsonbPathObjectsArrayKeyValueEqual(
    "metadata", 
    []string{"users", "preferences"}, 
    "theme", 
    "dark",
))
```

#### Combining Selectors

```go
// Apply multiple conditions
query.Apply(
    bunhelpers.WhereEqual("status", "active"),
    bunhelpers.WhereNotNull("verified_at"),
)

// Conditional application
isAdmin := true
query.Apply(bunhelpers.ApplyIf(isAdmin, 
    bunhelpers.WhereEqual("role", "admin"),
))

// OR groups
query.Apply(bunhelpers.OrGroup(
    bunhelpers.WhereEqual("status", "active"),
    bunhelpers.WhereEqual("status", "pending"),
))

// AND groups
query.Apply(bunhelpers.AndGroup(
    bunhelpers.WhereEqual("verified", true),
    bunhelpers.WhereNotNull("email"),
))

// Complex OR conditions
query.Apply(bunhelpers.Or(
    bunhelpers.WhereEqual("role", "admin"),
    bunhelpers.WhereEqual("role", "moderator"),
    bunhelpers.WhereEqual("role", "editor"),
))
```

### 2. Transaction Management

#### Simple Transactions with InTx

The `InTx` function handles transaction lifecycle automatically:

```go
err := bunhelpers.InTx(ctx, db, func(ctx context.Context) error {
    // Create user
    _, err := db.NewInsert().
        Model(&user).
        Exec(ctx)  // Automatically uses transaction
    if err != nil {
        return err // Transaction will rollback
    }
    
    // Create profile
    _, err = db.NewInsert().
        Model(&profile).
        Exec(ctx)
    if err != nil {
        return err // Transaction will rollback
    }
    
    return nil // Transaction will commit
})
```

#### Nested Transactions

`InTx` supports nested calls - only the outermost call creates the transaction:

```go
func CreateUser(ctx context.Context, db *bun.DB, user *User) error {
    return bunhelpers.InTx(ctx, db, func(ctx context.Context) error {
        // This might create a new transaction OR use existing one
        _, err := db.NewInsert().Model(user).Exec(ctx)
        return err
    })
}

func CreateUserWithProfile(ctx context.Context, db *bun.DB, user *User, profile *Profile) error {
    return bunhelpers.InTx(ctx, db, func(ctx context.Context) error {
        // Outer transaction
        if err := CreateUser(ctx, db, user); err != nil {
            return err // Nested call uses same transaction
        }
        
        _, err := db.NewInsert().Model(profile).Exec(ctx)
        return err
    })
}
```

#### Manual Transaction Management

For more control, use `TxToContext` and `TxFromContext`:

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()

// Store transaction in context
ctx = bunhelpers.TxToContext(ctx, &tx)

// Pass context to functions
err = createUser(ctx, db)
if err != nil {
    return err
}

return tx.Commit()
```

### 3. Querier Interface

The Querier interface provides context-aware query builders:

```go
type UserRepository struct {
    querier bunhelpers.Querier
}

func NewUserRepository(db *bun.DB) *UserRepository {
    return &UserRepository{
        querier: bunhelpers.NewQuerier(db),
    }
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    user := new(User)
    
    // Automatically uses transaction from context if available
    err := r.querier.NewSelectQuery(ctx).
        Model(user).
        Where("id = ?", id).
        Scan(ctx)
    
    if bunhelpers.IsNotFoundError(err) {
        return nil, fmt.Errorf("user not found")
    }
    
    return user, err
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    _, err := r.querier.NewInsertQuery(ctx).
        Model(user).
        Exec(ctx)
    
    if bunhelpers.IsConstraintError(err) {
        return fmt.Errorf("user already exists")
    }
    
    return err
}
```

### 4. Where Struct - Advanced Filtering

Use the pre-built `Where` struct for common filtering patterns:

```go
// Build complex filters
where := &bunhelpers.Where{
    IDs: []string{"1", "2", "3"},
    NotInIDs: []string{"99"},
    OnlyDeleted: false,
    WithDeleted: false,
    Limit: ToPtr(10),
    Offset: ToPtr(0),
    CreatedAfter: ToPtr(time.Now().Add(-24*time.Hour).UnixMilli()),
    SelectColumns: []string{"id", "name", "email"},
    SortBy: 1,
    SortDesc: true,
}

// Define sort mapping
where.Order = bunhelpers.Order{
    1: "created_at",
    2: "updated_at",
    3: "name",
}

// Apply to query
query := db.NewSelect().Model(&users)
query = where.Where(query)   // Apply WHERE conditions
query = where.Select(query)  // Apply SELECT, LIMIT, ORDER BY
```

#### Bitwise Flag Filtering

```go
const (
    FlagActive   = 1 << 0  // 1
    FlagVerified = 1 << 1  // 2
    FlagPremium  = 1 << 2  // 4
)

where := &bunhelpers.Where{
    HasFlags: []int{FlagActive, FlagVerified},    // Must have both flags
    HasNotFlags: []int{FlagPremium},              // Must not have flag
}

query := db.NewSelect().Model(&users)
query = where.Where(query)
```

#### Use as Selector

```go
where := &bunhelpers.Where{IDs: []string{"1", "2"}}

query := db.NewSelect().
    Model(&users).
    Apply(bunhelpers.NestedWhere(where)).
    Apply(bunhelpers.WhereEqual("status", "active"))
```

### 5. Error Handling

Specialized error handlers for common database scenarios:

```go
user := new(User)
err := db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)

if bunhelpers.IsNotFoundError(err) {
    return fmt.Errorf("user not found")
}

// Check constraint violations
_, err = db.NewInsert().Model(user).Exec(ctx)

if bunhelpers.IsConstraintError(err) {
    return fmt.Errorf("user with this email already exists")
}
```

#### Predefined Errors

```go
var (
    ErrInvalidRequest  = errors.New("invalid request")
    ErrInvalidResponse = errors.New("invalid response type")
    ErrEmptyPK         = errors.New("primary key (id, email etc.) is empty")
    ErrNilRequest      = errors.New("no request data")
)
```

## Complete Example

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"

    "github.com/uagolang/bunhelpers"
    "github.com/uptrace/bun"
    "github.com/uptrace/bun/dialect/pgdialect"
    "github.com/uptrace/bun/driver/pgdriver"
)

type User struct {
    bun.BaseModel `bun:"table:users"`
    
    ID        string    `bun:"id,pk"`
    Email     string    `bun:"email,notnull,unique"`
    Name      string    `bun:"name"`
    Status    string    `bun:"status"`
    Flags     int       `bun:"flags"`
    CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

type UserRepository struct {
    querier bunhelpers.Querier
}

func NewUserRepository(db *bun.DB) *UserRepository {
    return &UserRepository{
        querier: bunhelpers.NewQuerier(db),
    }
}

func (r *UserRepository) FindActive(ctx context.Context, limit int) ([]*User, error) {
    var users []*User
    
    err := r.querier.NewSelectQuery(ctx).
        Model(&users).
        Apply(bunhelpers.WhereEqual("status", "active")).
        Limit(limit).
        Order("created_at DESC").
        Scan(ctx)
    
    return users, err
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    _, err := r.querier.NewInsertQuery(ctx).
        Model(user).
        Exec(ctx)
    
    if bunhelpers.IsConstraintError(err) {
        return fmt.Errorf("user already exists")
    }
    
    return err
}

func main() {
    // Setup database
    dsn := "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
    sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
    db := bun.NewDB(sqldb, pgdialect.New())
    defer db.Close()
    
    repo := NewUserRepository(db)
    ctx := context.Background()
    
    // Example 1: Simple query with selectors
    users, err := repo.FindActive(ctx, 10)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d active users\n", len(users))
    
    // Example 2: Transaction with InTx
    err = bunhelpers.InTx(ctx, db, func(ctx context.Context) error {
        user := &User{
            ID:     "123",
            Email:  "test@example.com",
            Name:   "Test User",
            Status: "active",
        }
        
        return repo.Create(ctx, user)
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Example 3: Using Where struct
    where := &bunhelpers.Where{
        IDs: []string{"1", "2", "3"},
        Limit: ToPtr(10),
    }
    where.Order = bunhelpers.Order{1: "created_at"}
    where.SortBy = 1
    where.SortDesc = true
    
    query := db.NewSelect().Model(&users)
    query = where.Where(query)
    query = where.Select(query)
    
    if err := query.Scan(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## API Reference

### Selectors

- `Apply(selectors ...Selector) Selector` - Combine multiple selectors
- `ApplyIf(cond bool, selectors ...Selector) Selector` - Conditionally apply selectors
- `OrGroup(selectors ...Selector) Selector` - Create OR group
- `AndGroup(selectors ...Selector) Selector` - Create AND group
- `Or(selectors ...Selector) Selector` - Separate conditions with OR
- `NestedWhere(where *Where) Selector` - Use Where struct as selector
- `WhereEqual(col string, value any) Selector`
- `WhereNotEqual(col string, value any) Selector`
- `WhereNull(col string) Selector`
- `WhereNotNull(col string) Selector`
- `WhereIn(col string, values any) Selector`
- `WhereNotIn(col string, values any) Selector`
- `WhereContains(col string, substr string) Selector`
- `WhereBegins(col string, substr string) Selector`
- `WhereEnds(col string, substr string) Selector`
- `WhereBefore(col string, t time.Time) Selector`
- `WhereAfter(col string, t time.Time) Selector`
- `WhereDistinctOn(col string) Selector`
- `WhereJsonbEqual(col string, field string, value any) Selector`
- `WhereJsonbPathEqual(col string, path []string, value any) Selector`
- `WhereJsonbObjectsArrayKeyValueEqual(col string, key, field string, value any) Selector`
- `WhereJsonbPathObjectsArrayKeyValueEqual(col string, path []string, field string, value any) Selector`

### Transaction Context

- `InTx(ctx context.Context, client *bun.DB, fn func(ctx context.Context) error) error` - Execute function in transaction
- `TxToContext(ctx context.Context, tx *bun.Tx) context.Context` - Store transaction in context
- `TxFromContext(ctx context.Context) *bun.Tx` - Retrieve transaction from context

### Error Handling

- `IsNotFoundError(err error) bool` - Check if error is `sql.ErrNoRows`
- `IsConstraintError(err error) bool` - Check for unique constraint violations

### Querier Interface

- `NewQuerier(c *bun.DB) Querier` - Create new querier
- `NewSelectQuery(ctx context.Context) *bun.SelectQuery` - Get context-aware SELECT query
- `NewInsertQuery(ctx context.Context) *bun.InsertQuery` - Get context-aware INSERT query
- `NewUpdateQuery(ctx context.Context) *bun.UpdateQuery` - Get context-aware UPDATE query
- `NewDeleteQuery(ctx context.Context) *bun.DeleteQuery` - Get context-aware DELETE query

### Utilities

- `OrderAsc(col string) string` - Create ascending order expression
- `OrderDesc(col string) string` - Create descending order expression

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built on top of the excellent [Bun ORM](https://github.com/uptrace/bun)
