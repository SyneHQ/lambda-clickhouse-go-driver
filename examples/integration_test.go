//go:build integration
// +build integration

package clickhouse

import (
	"database/sql"
	"testing"
)

// Run with: go test -tags=integration -v
func TestIntegration(t *testing.T) {
	// Use your actual Lambda function details
	dsn := "clickhouse-lambda://clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/test.csv"

	db, err := sql.Open("clickhouse-lambda", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Simple query test
	query := "SELECT * FROM table LIMIT 2"
	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("Failed to get columns: %v", err)
	}

	t.Logf("Columns: %v", columns)

	rowCount := 0
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			t.Errorf("Scan failed: %v", err)
			continue
		}

		t.Logf("Row %d: %v", rowCount+1, values)
		rowCount++
	}

	if err := rows.Err(); err != nil {
		t.Errorf("Rows error: %v", err)
	}

	if rowCount == 0 {
		t.Error("Expected at least one row, got none")
	}

	t.Logf("Successfully processed %d rows", rowCount)
}
