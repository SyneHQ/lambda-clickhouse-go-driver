// Package clickhouse provides a database/sql driver for ClickHouse running on AWS Lambda
package clickhouse

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awscredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

const DriverName = "clickhouse-lambda"

// Driver implements the database/sql/driver.Driver interface
type Driver struct{}

// Config holds the configuration for the ClickHouse Lambda connection
type Config struct {
	FunctionName    string
	Region          string
	BucketName      string
	BucketPath      string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string // Optional for temporary credentials
}

// Open parses the data source name and returns a connection
// DSN format: clickhouse-lambda://function-name@region/bucket-name/path?param=value
func (d *Driver) Open(dsn string) (driver.Conn, error) {
	config, err := parseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Load AWS config with credentials if provided
	var configOptions []func(*awsconfig.LoadOptions) error
	configOptions = append(configOptions, awsconfig.WithRegion(config.Region))

	// Use provided credentials if available
	if config.AccessKeyID != "" && config.SecretAccessKey != "" {
		credentials := awscredentials.NewStaticCredentialsProvider(
			config.AccessKeyID,
			config.SecretAccessKey,
			config.SessionToken,
		)
		configOptions = append(configOptions, awsconfig.WithCredentialsProvider(credentials))
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), configOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	return &Conn{
		config:       config,
		lambdaClient: lambdaClient,
	}, nil
}

// parseDSN parses the data source name
// Format: clickhouse-lambda://function-name@region/bucket-name/path?aws_access_key_id=xxx&aws_secret_access_key=yyy&aws_session_token=zzz
func parseDSN(dsn string) (*Config, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "clickhouse-lambda" {
		return nil, fmt.Errorf("invalid scheme: %s", u.Scheme)
	}

	// Extract function name from user info
	functionName := u.User.Username()
	if functionName == "" {
		return nil, fmt.Errorf("function name is required")
	}

	// Extract region from host
	region := u.Host
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}

	// Extract bucket name and path
	pathParts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(pathParts) < 1 {
		return nil, fmt.Errorf("bucket name is required")
	}

	bucketName := pathParts[0]
	bucketPath := "/"
	if len(pathParts) > 1 {
		// Join the remaining path parts and URL decode
		rawPath := strings.Join(pathParts[1:], "/")
		decodedPath, err := url.QueryUnescape(rawPath)
		if err != nil {
			// If decoding fails, use the raw path
			bucketPath = "/" + rawPath
		} else {
			bucketPath = "/" + decodedPath
		}
	}

	// Extract AWS credentials from query parameters
	queryParams := u.Query()
	accessKeyID := queryParams.Get("aws_access_key_id")
	secretAccessKey := queryParams.Get("aws_secret_access_key")
	sessionToken := queryParams.Get("aws_session_token")

	return &Config{
		FunctionName:    functionName,
		Region:          region,
		BucketName:      bucketName,
		BucketPath:      bucketPath,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		SessionToken:    sessionToken,
	}, nil
}

func init() {
	sql.Register(DriverName, &Driver{})
}
