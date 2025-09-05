# ClickHouse Lambda Driver Documentation

## Overview

The ClickHouse Lambda Driver is a Go database/sql driver that enables you to query ClickHouse running on AWS Lambda. It supports both S3 and HTTPS data sources, making it easy to work with both your own data and public datasets.

## Features

- **Dual Data Source Support**: Query both S3 files and HTTPS URLs
- **Standard Interface**: Implements the standard `database/sql` interface
- **AWS Integration**: Seamless integration with AWS Lambda and S3
- **Public Dataset Access**: Direct access to public datasets via HTTPS URLs
- **Flexible Authentication**: Supports IAM roles, environment variables, and explicit credentials

## Quick Start

### Installation

```bash
go get github.com/synehq/lambda-clickhouse-go-driver/pkg/clickhouse
```

### Basic Usage

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/synehq/lambda-clickhouse-go-driver/pkg/clickhouse"
)

func main() {
    // Connect to ClickHouse Lambda
    dsn := "clickhouse-lambda://function-name@region/bucket-name/path"
    db, err := sql.Open("clickhouse-lambda", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Query your data
    rows, err := db.Query("SELECT * FROM table LIMIT 5")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    // Process results...
}
```

## Data Source Configuration

### S3 Data Sources

For S3 files, use the standard bucket/path format:

```go
dsn := "clickhouse-lambda://function@region/bucket/file.csv"
```

### HTTPS Data Sources

For HTTPS URLs, URL-encode the full URL in the path:

```go
dsn := "clickhouse-lambda://function@region/bucket/https%3A%2F%2Fexample.com%2Fdata.csv"
```

Or use domain-only paths for common data sources:

```go
dsn := "clickhouse-lambda://function@region/bucket/raw.githubusercontent.com/user/repo/data.csv"
```

## Authentication

The driver supports multiple authentication methods:

### 1. IAM Roles (Recommended)
```go
// No credentials needed - uses IAM role
dsn := "clickhouse-lambda://function@region/bucket/path"
```

### 2. Environment Variables
```bash
# Standard AWS credentials
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key

# or use driver-specific environment variables
export LAMBDA_CLICKHOUSE_ROLE_AWS_SECRET_ACCESS_KEY=your_secret_key
export LAMBDA_CLICKHOUSE_ROLE_AWS_ACCESS_KEY_ID=your_access_key

# Set the AWS region
export AWS_REGION=us-east-1
```

### 3. Explicit Credentials
```go
dsn := "clickhouse-lambda://function@region/bucket/path?aws_access_key_id=KEY&aws_secret_access_key=SECRET"
```

## Supported Operations

### CREATE TABLE
```go
createQuery := `
    CREATE TABLE users (
        id UInt32,
        name String, 
        email String,
        age UInt8
    ) ENGINE = Memory
`
_, err := db.Exec(createQuery)
```

### SELECT Queries
```go
rows, err := db.Query("SELECT * FROM table WHERE age > 25 LIMIT 10")
```

### Prepared Statements
```go
stmt, err := db.Prepare("SELECT count(*) FROM table WHERE age > ?")
var count int64
err = stmt.QueryRow(25).Scan(&count)
```

### Aggregations
```go
var avgAge float64
err = db.QueryRow("SELECT avg(age) FROM table").Scan(&avgAge)
```

## Error Handling

The driver provides detailed error information:

```go
rows, err := db.Query("SELECT * FROM table")
if err != nil {
    // Check for specific error types
    if strings.Contains(err.Error(), "AccessDenied") {
        log.Println("AWS credentials issue")
    } else if strings.Contains(err.Error(), "Function not found") {
        log.Println("Lambda function not found")
    } else {
        log.Printf("Query error: %v", err)
    }
}
```

## Performance Considerations

- **Cold Starts**: Lambda functions may have cold start delays
- **Memory Usage**: Large datasets may require increased Lambda memory
- **Timeout**: Set appropriate timeouts for long-running queries
- **Concurrency**: Consider Lambda concurrency limits

## Troubleshooting

### Common Issues

1. **Access Denied**: Check AWS credentials and permissions
2. **Function Not Found**: Verify Lambda function name and region
3. **Timeout**: Increase Lambda timeout or optimize queries
4. **Memory Issues**: Increase Lambda memory allocation

### Debug Mode

Enable debug logging by setting the log level:

```go
// In your Lambda function
console.log("Debug: Query executed", query);
```

## Examples

See the [examples/](../examples/) directory for complete working examples:

- `comprehensive_example.go` - Complete feature demonstration
- `test_enhanced_https.go` - HTTPS URL support testing
- `advanced_example.go` - Advanced usage patterns

## API Reference

### Driver Interface

The driver implements the standard `database/sql/driver` interface:

```go
type Driver interface {
    Open(name string) (Conn, error)
}
```

### Connection Interface

```go
type Conn interface {
    Prepare(query string) (Stmt, error)
    Close() error
    Begin() (Tx, error)
}
```

### Statement Interface

```go
type Stmt interface {
    Close() error
    NumInput() int
    Exec(args []Value) (Result, error)
    Query(args []Value) (Rows, error)
}
```

## License

MIT License - see LICENSE file for details.