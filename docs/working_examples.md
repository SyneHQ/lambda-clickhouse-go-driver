# ClickHouse Lambda - What Actually Works

Based on the test results, here's what your ClickHouse Lambda driver **can and cannot** do:

## ✅ **FULLY SUPPORTED OPERATIONS**

### 1. CREATE TABLE with Schema ✅
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
// ✅ This WORKS!
```

### 2. SELECT Queries on S3 Data ✅
```go
// The 'table' is automatically created from your S3 CSV file
rows, err := db.Query("SELECT * FROM table WHERE _4 > 25 LIMIT 5")
// ✅ This WORKS! (_1, _2, _3, _4... are column names)
```

### 3. Complex Analytics on S3 Data ✅
```go
// Aggregation queries
db.QueryRow("SELECT avg(_4), count(*) FROM table").Scan(&avgAge, &count)

// GROUP BY queries  
rows, err := db.Query("SELECT _7, count(*) FROM table GROUP BY _7 ORDER BY count(*) DESC")

// WHERE clauses with complex conditions
rows, err := db.Query("SELECT * FROM table WHERE _4 > 30 AND _7 = 'CA'")
// ✅ All of these WORK!
```

### 4. Loading from Different CSV Files ✅
```go
// Load from different CSV files by changing the DSN path:

// Load from users.csv
dsn1 := "clickhouse-lambda://function@region/bucket/users.csv"

// Load from sales.csv
dsn2 := "clickhouse-lambda://function@region/bucket/sales/2024/january.csv" 

// Load from nested paths
dsn3 := "clickhouse-lambda://function@region/bucket/data/transactions/daily/2024-01-15.csv"
// ✅ All of these WORK!
```

## ⚠️ **LIMITATIONS** 

### 1. CREATE TABLE AS SELECT Needs ENGINE
```go
// ❌ This FAILS:
createAsSelectQuery := `CREATE TABLE people AS SELECT _1 as id FROM table`

// ✅ This WORKS:
createAsSelectQuery := `CREATE TABLE people ENGINE = Memory AS SELECT _1 as id FROM table`
```

### 2. Tables Don't Persist Between Lambda Invocations
- Each Lambda invocation starts fresh
- Created tables only exist during that single query execution
- This is **expected behavior** for serverless ClickHouse

### 3. Multi-Statement Queries Work Within Single Invocation
```go
// ✅ This WORKS - multiple statements in one query:
query := `
    CREATE TABLE temp_stats ENGINE = Memory AS 
    SELECT _7 as state, avg(_4) as avg_age FROM table GROUP BY _7;
    
    SELECT * FROM temp_stats WHERE avg_age > 30;
`
rows, err := db.Query(query)
```

## 🚀 **RECOMMENDED USAGE PATTERNS**

### Pattern 1: Direct S3 Analytics
```go
// Best for: One-off analytics queries
dsn := "clickhouse-lambda://function@region/bucket/sales_2024.csv"
db, _ := sql.Open("clickhouse-lambda", dsn)

// Direct queries on the CSV data
avgSales, _ := db.QueryRow("SELECT avg(_5) FROM table WHERE _2 = 'January'")
```

### Pattern 2: Complex Analytics with Temp Tables
```go
// Best for: Complex multi-step analytics
query := `
    CREATE TABLE sales_summary ENGINE = Memory AS 
    SELECT _2 as month, _3 as region, sum(_5) as total_sales 
    FROM table 
    GROUP BY _2, _3;
    
    SELECT region, avg(total_sales) as avg_monthly_sales 
    FROM sales_summary 
    GROUP BY region 
    ORDER BY avg_monthly_sales DESC;
`
```

### Pattern 3: Multi-File Analysis
```go
// Load different datasets by changing DSN
datasets := []string{
    "users.csv",
    "sales/2024/q1.csv", 
    "analytics/metrics.csv",
}

for _, dataset := range datasets {
    dsn := fmt.Sprintf("clickhouse-lambda://function@region/bucket/%s", dataset)
    db, _ := sql.Open("clickhouse-lambda", dsn)
    
    // Analyze each dataset
    rows, _ := db.Query("SELECT count(*), avg(_1) FROM table")
    // Process results...
}
```

## 📊 **Column Naming in CSV Files**

ClickHouse automatically names CSV columns as `_1`, `_2`, `_3`, etc. You can alias them:

```go
query := `
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
    WHERE age > 25
`
```

## 🎯 **Perfect Use Cases**

1. **ETL Processing**: Load CSV → Transform → Output results
2. **Analytics Dashboards**: Run complex queries on S3 data lakes  
3. **Data Validation**: Check data quality across multiple CSV files
4. **Reporting**: Generate aggregated reports from raw CSV data
5. **Data Migration**: Transform data formats using ClickHouse SQL
