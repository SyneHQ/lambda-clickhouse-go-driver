//go:build integration
// +build integration

package clickhouse

import (
	"database/sql"
	"testing"
)

// Run with: go test -tags=integration -v -run TestAdvancedOperations
func TestAdvancedOperations(t *testing.T) {
	dsn := "clickhouse-lambda://ClickhouseLambdaStack-ClickhouseLambda3697569D-iktHR8bNJ4BS@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/test.csv"

	db, err := sql.Open("clickhouse-lambda", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test 1: CREATE TABLE with explicit schema
	t.Log("Testing CREATE TABLE with explicit schema")
	createQuery := `
		CREATE TABLE users (
			id UInt32,
			name String, 
			email String,
			age UInt8
		) ENGINE = Memory
	`

	_, err = db.Exec(createQuery)
	if err != nil {
		t.Logf("CREATE TABLE failed (expected if table exists): %v", err)
	} else {
		t.Log("✓ CREATE TABLE succeeded")
	}

	// Test 2: CREATE TABLE AS SELECT from S3 data
	t.Log("Testing CREATE TABLE AS SELECT from S3 data")
	createAsSelectQuery := `
		CREATE TABLE people AS 
		SELECT 
			_1 as id,
			_2 as first_name, 
			_3 as last_name,
			_4 as age,
			_5 as street,
			_6 as city,
			_7 as state,
			_8 as zip,
			_9 as income
		FROM table
	`

	_, err = db.Exec(createAsSelectQuery)
	if err != nil {
		t.Logf("CREATE TABLE AS SELECT failed (expected if table exists): %v", err)
	} else {
		t.Log("✓ CREATE TABLE AS SELECT succeeded")
	}

	// Test 3: Query created table with WHERE clause
	t.Log("Testing query with WHERE clause")
	rows, err := db.Query("SELECT first_name, last_name, age FROM people WHERE age > 25 LIMIT 3")
	if err != nil {
		t.Errorf("WHERE query failed: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		var firstName, lastName string
		var age int64
		if err := rows.Scan(&firstName, &lastName, &age); err != nil {
			t.Errorf("Scan failed: %v", err)
			continue
		}
		t.Logf("  Row %d: %s %s (age %d)", rowCount+1, firstName, lastName, age)
		rowCount++
	}

	if rowCount == 0 {
		t.Error("Expected at least one row from WHERE query")
	} else {
		t.Logf("✓ WHERE query returned %d rows", rowCount)
	}

	// Test 4: Aggregation query
	t.Log("Testing aggregation query")
	var avgAge float64
	var count int64
	err = db.QueryRow("SELECT avg(age), count(*) FROM people").Scan(&avgAge, &count)
	if err != nil {
		t.Errorf("Aggregation query failed: %v", err)
	} else {
		t.Logf("✓ Aggregation: Average age %.1f, Total count %d", avgAge, count)
	}

	// Test 5: GROUP BY query
	t.Log("Testing GROUP BY query")
	groupRows, err := db.Query("SELECT state, count(*) as cnt FROM people GROUP BY state ORDER BY cnt DESC LIMIT 3")
	if err != nil {
		t.Errorf("GROUP BY query failed: %v", err)
		return
	}
	defer groupRows.Close()

	groupCount := 0
	for groupRows.Next() {
		var state string
		var cnt int64
		if err := groupRows.Scan(&state, &cnt); err != nil {
			t.Errorf("GROUP BY scan failed: %v", err)
			continue
		}
		t.Logf("  State %s: %d people", state, cnt)
		groupCount++
	}

	if groupCount == 0 {
		t.Error("Expected at least one row from GROUP BY query")
	} else {
		t.Logf("✓ GROUP BY query returned %d states", groupCount)
	}
}
