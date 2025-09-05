package main

import (
	"database/sql"
	"fmt"
	"log"

	// Import the ClickHouse Lambda driver
	_ "github.com/synehq/lambda-clickhouse-go-driver/pkg/clickhouse"
)

func main() {
	// DSN format: clickhouse-lambda://function-name@region/bucket-name/path
	// Without credentials (uses default AWS credential chain):
	dsn := "clickhouse-lambda://clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/test.csv"

	// With AWS credentials in DSN (uncomment to use):
	// dsn := "clickhouse-lambda://clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/test.csv?aws_access_key_id=YOUR_ACCESS_KEY&aws_secret_access_key=YOUR_SECRET_KEY"

	// Open database connection
	db, err := sql.Open("clickhouse-lambda", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to ClickHouse Lambda successfully!")

	// Example query
	query := "SELECT * FROM table LIMIT 5"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to get columns: %v", err)
	}

	fmt.Printf("Columns: %v\n", columns)

	// Process results
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	fmt.Println("\nResults:")
	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			log.Printf("Scan failed: %v", err)
			continue
		}

		for i, col := range columns {
			fmt.Printf("%s: %v\t", col, values[i])
		}
		fmt.Println()
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Rows error: %v", err)
	}

	// Example with prepared statement
	fmt.Println("\nUsing prepared statement:")
	stmt, err := db.Prepare("SELECT count(*) FROM table")
	if err != nil {
		log.Fatalf("Prepare failed: %v", err)
	}
	defer stmt.Close()

	var count int64
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		log.Fatalf("QueryRow failed: %v", err)
	}

	fmt.Printf("Total rows in table: %d\n", count)
}
