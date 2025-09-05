# ClickHouse Lambda Driver

A Go database/sql driver for ClickHouse running on AWS Lambda with support for both S3 and HTTPS data sources.

## Features

- 🚀 **Dual Data Source Support**: Works with both S3 files and HTTPS URLs
- 🔗 **Standard database/sql Interface**: Drop-in replacement for other ClickHouse drivers
- ☁️ **AWS Integration**: Seamless integration with AWS Lambda and S3
- 📊 **Public Dataset Support**: Direct access to public datasets via HTTPS URLs
- 🔒 **AWS Credentials**: Supports IAM roles, environment variables, and explicit credentials

## Quick Start

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

## Data Source Support

### S3 Files
```go
dsn := "clickhouse-lambda://function@region/bucket/file.csv"
```

### HTTPS URLs
```go
dsn := "clickhouse-lambda://function@region/bucket/https%3A%2F%2Fexample.com%2Fdata.csv"
```

## Directory Structure

```
go-driver/
├── cmd/clickhouse-lambda/     # CLI application
├── pkg/clickhouse/            # Core driver package
├── examples/                  # Example applications and tests
├── docs/                      # Documentation
└── go.mod                     # Go module definition
```

## Documentation

- [Quick Start Guide](docs/QUICKSTART.md)
- [Working Examples](docs/working_examples.md)
- [API Reference](docs/README.md)

## Examples

See the [examples/](examples/) directory for complete working examples:

- `test_enhanced_https.go` - HTTPS URL support testing
- `test_public_csv.go` - Public CSV data access
- `advanced_example.go` - Advanced usage patterns

## Installation

```bash
go get github.com/synehq/lambda-clickhouse-go-driver/pkg/clickhouse
```

## License

MIT License - see LICENSE file for details.
