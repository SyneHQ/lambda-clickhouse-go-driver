package clickhouse

import (
	"context"
	"database/sql/driver"
	"fmt"
)

// Stmt implements the database/sql/driver.Stmt interface
type Stmt struct {
	conn  *Conn
	query string
}

// Close implements driver.Stmt.Close
func (s *Stmt) Close() error {
	return nil
}

// NumInput implements driver.Stmt.NumInput
func (s *Stmt) NumInput() int {
	// For now, we don't support parameterized queries
	return 0
}

// Exec implements driver.Stmt.Exec
func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("prepared statements with arguments not yet supported")
	}

	response, err := s.conn.Query(s.query)
	if err != nil {
		return nil, err
	}

	return &Result{
		response: response,
	}, nil
}

// Query implements driver.Stmt.Query
func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("prepared statements with arguments not yet supported")
	}

	response, err := s.conn.Query(s.query)
	if err != nil {
		return nil, err
	}

	return NewRows(response.Body), nil
}

// ExecContext implements driver.StmtExecContext
func (s *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("prepared statements with arguments not yet supported")
	}

	response, err := s.conn.Query(s.query)
	if err != nil {
		return nil, err
	}

	return &Result{
		response: response,
	}, nil
}

// QueryContext implements driver.StmtQueryContext
func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("prepared statements with arguments not yet supported")
	}

	response, err := s.conn.Query(s.query)
	if err != nil {
		return nil, err
	}

	return NewRows(response.Body), nil
}

// Result implements the database/sql/driver.Result interface
type Result struct {
	response *LambdaResponse
}

// LastInsertId implements driver.Result.LastInsertId
func (r *Result) LastInsertId() (int64, error) {
	// ClickHouse doesn't typically use auto-incrementing IDs
	return 0, fmt.Errorf("LastInsertId not supported")
}

// RowsAffected implements driver.Result.RowsAffected
func (r *Result) RowsAffected() (int64, error) {
	// For now, we can't determine rows affected from the response
	// This would need to be parsed from the ClickHouse response
	return 0, fmt.Errorf("RowsAffected not supported yet")
}
