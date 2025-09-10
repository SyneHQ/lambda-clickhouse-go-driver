package main

import (
	"database/sql"
	"fmt"

	_ "github.com/synehq/lambda-clickhouse-go-driver/pkg/clickhouse"
)

func main() {
	fmt.Println("=== Testing Enhanced HTTPS URL Support ===\n")

	// Test 1: GitHub CSV with protocol in DSN path
	fmt.Println("1. Testing GitHub CSV with full HTTPS URL:")
	dsn1 := "clickhouse-lambda://ClickhouseLambdaStack-ClickhouseLambda3697569D-iktHR8bNJ4BS@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/https%3A%2F%2Fraw.githubusercontent.com%2Fdatasciencedojo%2Fdatasets%2Fmaster%2Ftitanic.csv"

	db1, err := sql.Open("clickhouse-lambda", dsn1)
	if err != nil {
		fmt.Printf("   ❌ Failed to open: %v\n", err)
		return
	}
	defer db1.Close()

	var count1 int64
	err = db1.QueryRow("SELECT count(*) FROM table").Scan(&count1)
	if err != nil {
		fmt.Printf("   ❌ Query failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Count: %d rows\n", count1)
		if count1 == 891 {
			fmt.Println("   🎉 SUCCESS! This is the Titanic dataset (891 rows), not your S3 data!")
		} else if count1 == 100 {
			fmt.Println("   ⚠️ Still getting S3 data (100 rows) - enhancement may not be deployed yet")
		}
	}

	// Test 2: Direct domain path (no https prefix)
	fmt.Println("\n2. Testing domain path without https prefix:")
	dsn2 := "clickhouse-lambda://ClickhouseLambdaStack-ClickhouseLambda3697569D-iktHR8bNJ4BS@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/raw.githubusercontent.com/datasciencedojo/datasets/master/titanic.csv"

	db2, err := sql.Open("clickhouse-lambda", dsn2)
	if err != nil {
		fmt.Printf("   ❌ Failed to open: %v\n", err)
		return
	}
	defer db2.Close()

	var count2 int64
	err = db2.QueryRow("SELECT count(*) FROM table").Scan(&count2)
	if err != nil {
		fmt.Printf("   ❌ Query failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Count: %d rows\n", count2)
		if count2 == 891 {
			fmt.Println("   🎉 SUCCESS! Direct HTTPS URL detection is working!")
		}
	}

	// Test 3: Verify we can still use S3 data
	fmt.Println("\n3. Testing backward compatibility with S3:")
	dsn3 := "clickhouse-lambda://ClickhouseLambdaStack-ClickhouseLambda3697569D-iktHR8bNJ4BS@us-east-1/clickhouselambdastack-clickhousebucket38d56561-idhkkhidojc7/test.csv"

	db3, err := sql.Open("clickhouse-lambda", dsn3)
	if err != nil {
		fmt.Printf("   ❌ Failed to open: %v\n", err)
		return
	}
	defer db3.Close()

	var count3 int64
	err = db3.QueryRow("SELECT count(*) FROM table").Scan(&count3)
	if err != nil {
		fmt.Printf("   ❌ Query failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Count: %d rows (S3 data should still work)\n", count3)
	}

	fmt.Println("\n=== Test Complete ===")

	if count1 == 891 || count2 == 891 {
		fmt.Println("🚀 HTTPS URL support is working! You can now use public datasets directly in DSN!")
	} else {
		fmt.Println("📝 Note: The Lambda function may need to be redeployed with the enhanced code.")
		fmt.Println("💡 When deployed, you'll be able to use any public HTTPS CSV URL directly in your DSN!")
	}
}
