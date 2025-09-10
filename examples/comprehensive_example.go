package main

import (
	"database/sql"
	"fmt"

	_ "github.com/synehq/lambda-clickhouse-go-driver/pkg/clickhouse"
)

func main() {
	fmt.Println("=== ClickHouse Lambda Driver - Comprehensive Example ===\n")

	// Example 1: S3 Data Source
	fmt.Println("1. Testing S3 Data Source:")
	testS3DataSource()

	// Example 2: HTTPS Data Source
	fmt.Println("\n2. Testing HTTPS Data Source:")
	testHTTPSDataSource()

	// Example 3: Advanced Queries
	fmt.Println("\n3. Testing Advanced Queries:")
	testAdvancedQueries()
}

func testS3DataSource() {
	// Use your S3 bucket data
	dsn := "clickhouse-lambda://ClickhouseLambdaStack-ClickhouseLambda3697569D-iktHR8bNJ4BS@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/test.csv"

	db, err := sql.Open("clickhouse-lambda", dsn)
	if err != nil {
		fmt.Printf("   ❌ Failed to open S3 connection: %v\n", err)
		return
	}
	defer db.Close()

	var count int64
	err = db.QueryRow("SELECT count(*) FROM table").Scan(&count)
	if err != nil {
		fmt.Printf("   ❌ S3 query failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ S3 data count: %d rows\n", count)
	}
}

func testHTTPSDataSource() {
	// Use GitHub Titanic dataset
	dsn := "clickhouse-lambda://ClickhouseLambdaStack-ClickhouseLambda3697569D-iktHR8bNJ4BS@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/https%3A%2F%2Fraw.githubusercontent.com%2Fdatasciencedojo%2Fdatasets%2Fmaster%2Ftitanic.csv"

	db, err := sql.Open("clickhouse-lambda", dsn)
	if err != nil {
		fmt.Printf("   ❌ Failed to open HTTPS connection: %v\n", err)
		return
	}
	defer db.Close()

	var count int64
	err = db.QueryRow("SELECT count(*) FROM table").Scan(&count)
	if err != nil {
		fmt.Printf("   ❌ HTTPS query failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ HTTPS data count: %d rows (Titanic dataset)\n", count)
	}
}

func testAdvancedQueries() {
	// Test with HTTPS data for more interesting queries
	dsn := "clickhouse-lambda://ClickhouseLambdaStack-ClickhouseLambda3697569D-iktHR8bNJ4BS@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/https%3A%2F%2Fraw.githubusercontent.com%2Fdatasciencedojo%2Fdatasets%2Fmaster%2Ftitanic.csv"

	db, err := sql.Open("clickhouse-lambda", dsn)
	if err != nil {
		fmt.Printf("   ❌ Failed to open connection: %v\n", err)
		return
	}
	defer db.Close()

	// Test basic data access
	rows, err := db.Query("SELECT * FROM table LIMIT 3")
	if err != nil {
		fmt.Printf("   ❌ SELECT query failed: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("   📊 Sample data (first 3 rows):")
	columns, _ := rows.Columns()
	fmt.Printf("      Columns: %v\n", columns)

	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	rowCount := 0
	for rows.Next() && rowCount < 3 {
		if err := rows.Scan(scanArgs...); err != nil {
			continue
		}
		fmt.Printf("      Row %d: ", rowCount+1)
		for i, col := range columns {
			fmt.Printf("%s=%v ", col, values[i])
		}
		fmt.Println()
		rowCount++
	}

	// Test simple aggregation
	var totalCount int64
	err = db.QueryRow("SELECT count(*) FROM table").Scan(&totalCount)
	if err != nil {
		fmt.Printf("   ❌ Count query failed: %v\n", err)
	} else {
		fmt.Printf("   📈 Total records: %d\n", totalCount)
	}
}
