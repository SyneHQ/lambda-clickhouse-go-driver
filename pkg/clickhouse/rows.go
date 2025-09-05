package clickhouse

import (
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Rows implements the database/sql/driver.Rows interface
type Rows struct {
	data    [][]string
	columns []string
	pos     int
	closed  bool
}

// NewRows creates a new Rows instance from tab-separated data
func NewRows(data string) *Rows {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) == 0 {
		return &Rows{
			data:    [][]string{},
			columns: []string{},
			pos:     0,
		}
	}

	// Parse the data - assuming tab-separated values
	var rows [][]string
	var columns []string

	for i, line := range lines {
		if line == "" {
			continue
		}
		
		fields := strings.Split(line, "\t")
		
		// For the first row, we'll infer column names
		if i == 0 && len(columns) == 0 {
			// Generate column names based on the number of fields
			for j := range fields {
				columns = append(columns, fmt.Sprintf("col%d", j+1))
			}
		}
		
		rows = append(rows, fields)
	}

	return &Rows{
		data:    rows,
		columns: columns,
		pos:     -1, // Start before the first row
	}
}

// Columns implements driver.Rows.Columns
func (r *Rows) Columns() []string {
	return r.columns
}

// Close implements driver.Rows.Close
func (r *Rows) Close() error {
	r.closed = true
	return nil
}

// Next implements driver.Rows.Next
func (r *Rows) Next(dest []driver.Value) error {
	if r.closed {
		return io.EOF
	}

	r.pos++
	if r.pos >= len(r.data) {
		return io.EOF
	}

	row := r.data[r.pos]
	if len(dest) != len(row) {
		return fmt.Errorf("expected %d columns, got %d", len(row), len(dest))
	}

	// Convert string values to appropriate types
	for i, val := range row {
		dest[i] = convertValue(val)
	}

	return nil
}

// convertValue attempts to convert a string value to the most appropriate Go type
func convertValue(s string) driver.Value {
	if s == "" {
		return nil
	}

	// Try to parse as integer
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}

	// Try to parse as float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	// Try to parse as boolean
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}

	// Default to string
	return s
}

// ColumnTypeScanType implements driver.RowsColumnTypeScanType (optional interface)
func (r *Rows) ColumnTypeScanType(index int) reflect.Type {
	// For now, return interface{} type - could be improved with type inference
	return reflect.TypeOf((*interface{})(nil)).Elem()
}

// HasNextResultSet implements driver.RowsNextResultSet (optional interface)
func (r *Rows) HasNextResultSet() bool {
	return false
}

// NextResultSet implements driver.RowsNextResultSet (optional interface)  
func (r *Rows) NextResultSet() error {
	return io.EOF
}
