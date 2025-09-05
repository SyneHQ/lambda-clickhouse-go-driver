package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "clickhouse-lambda-driver/pkg/clickhouse"
)

func main() {
	// Connect to ClickHouse Lambda
	dsn := "clickhouse-lambda://clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/test.csv"

	db, err := sql.Open("clickhouse-lambda", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("=== Advanced ClickHouse Lambda Operations ===\n")

	// 1. CREATE TABLE with specific schema
	fmt.Println("1. Creating a table with explicit schema:")
	createQuery := `
		CREATE TABLE users (
			id UInt32,
			name String,
			email String,
			age UInt8
		) ENGINE = Memory
	`

	if _, err := db.Exec(createQuery); err != nil {
		log.Printf("Create table failed: %v", err)
	} else {
		fmt.Println("✓ Table 'users' created successfully")
	}

	// 2. CREATE TABLE AS SELECT (from the default S3 data)
	fmt.Println("\n2. Creating table from S3 data with column names:")
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

	if _, err := db.Exec(createAsSelectQuery); err != nil {
		log.Printf("Create table as select failed: %v", err)
	} else {
		fmt.Println("✓ Table 'people' created from S3 data")
	}

	// 3. Query the created table
	fmt.Println("\n3. Querying the created table:")
	rows, err := db.Query("SELECT first_name, last_name, age FROM people WHERE age > 30 LIMIT 3")
	if err != nil {
		log.Printf("Query failed: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var firstName, lastName string
			var age int
			if err := rows.Scan(&firstName, &lastName, &age); err != nil {
				log.Printf("Scan failed: %v", err)
				continue
			}
			fmt.Printf("  %s %s (age %d)\n", firstName, lastName, age)
		}
	}

	// 4. Aggregation queries
	fmt.Println("\n4. Aggregation example:")
	var avgAge float64
	var count int
	err = db.QueryRow("SELECT avg(age), count(*) FROM people").Scan(&avgAge, &count)
	if err != nil {
		log.Printf("Aggregation query failed: %v", err)
	} else {
		fmt.Printf("  Average age: %.1f, Total people: %d\n", avgAge, count)
	}

	// 5. Group by operations
	fmt.Println("\n5. Group by state:")
	stateRows, err := db.Query("SELECT state, count(*) as count FROM people GROUP BY state ORDER BY count DESC LIMIT 5")
	if err != nil {
		log.Printf("Group by query failed: %v", err)
	} else {
		defer stateRows.Close()
		for stateRows.Next() {
			var state string
			var count int
			if err := stateRows.Scan(&state, &count); err != nil {
				log.Printf("Scan failed: %v", err)
				continue
			}
			fmt.Printf("  %s: %d people\n", state, count)
		}
	}

	fmt.Println("\n=== Advanced Operations Complete ===")
}

// Example function to demonstrate loading from different CSV files
func loadFromDifferentCSV() {
	// To load from a different CSV file, change the path in the DSN
	// For example, if you have multiple CSV files in your bucket:

	// Load from users.csv
	usersDSN := "clickhouse-lambda://clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/users.csv"

	// Load from sales.csv
	salesDSN := "clickhouse-lambda://clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/sales.csv"

	// Load from nested path
	nestedDSN := "clickhouse-lambda://clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/data/2024/january/transactions.csv"

	fmt.Printf("Example DSNs for different CSV files:\n")
	fmt.Printf("Users: %s\n", usersDSN)
	fmt.Printf("Sales: %s\n", salesDSN)
	fmt.Printf("Nested: %s\n", nestedDSN)
}
