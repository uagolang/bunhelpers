package bunhelpers

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// newTestDB creates a test database with a mock driver
func newTestDB() *bun.DB {
	sqldb := sql.OpenDB(&mockConnector{})
	return bun.NewDB(sqldb, pgdialect.New())
}

// Mock driver implementation for testing
type mockConnector struct{}

func (c *mockConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return &mockConn{}, nil
}

func (c *mockConnector) Driver() driver.Driver {
	return &mockDriver{}
}

type mockDriver struct{}

func (d *mockDriver) Open(name string) (driver.Conn, error) {
	return &mockConn{}, nil
}

type mockConn struct{}

func (c *mockConn) Prepare(query string) (driver.Stmt, error) {
	return &mockStmt{}, nil
}

func (c *mockConn) Close() error {
	return nil
}

func (c *mockConn) Begin() (driver.Tx, error) {
	return &mockTx{}, nil
}

type mockStmt struct{}

func (s *mockStmt) Close() error {
	return nil
}

func (s *mockStmt) NumInput() int {
	return 0
}

func (s *mockStmt) Exec(args []driver.Value) (driver.Result, error) {
	return &mockResult{}, nil
}

func (s *mockStmt) Query(args []driver.Value) (driver.Rows, error) {
	return &mockRows{}, nil
}

type mockTx struct{}

func (tx *mockTx) Commit() error {
	return nil
}

func (tx *mockTx) Rollback() error {
	return nil
}

type mockResult struct{}

func (r *mockResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (r *mockResult) RowsAffected() (int64, error) {
	return 0, nil
}

type mockRows struct{}

func (r *mockRows) Columns() []string {
	return []string{}
}

func (r *mockRows) Close() error {
	return nil
}

func (r *mockRows) Next(dest []driver.Value) error {
	return sql.ErrNoRows
}
