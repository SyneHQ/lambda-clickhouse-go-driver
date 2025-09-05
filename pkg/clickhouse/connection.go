package clickhouse

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// Conn implements the database/sql/driver.Conn interface
type Conn struct {
	config       *Config
	lambdaClient *lambda.Client
	closed       bool
}

// LambdaEvent represents the event structure expected by the Lambda function
type LambdaEvent struct {
	RawPath        string                 `json:"rawPath"`
	RequestContext RequestContext         `json:"requestContext"`
	Body           string                 `json:"body"`
}

// RequestContext represents the request context for the Lambda event
type RequestContext struct {
	HTTP HTTPContext `json:"http"`
}

// HTTPContext represents the HTTP context
type HTTPContext struct {
	Method string `json:"method"`
}

// LambdaResponse represents the response from the Lambda function
type LambdaResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

// Prepare implements driver.Conn.Prepare
func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	if c.closed {
		return nil, driver.ErrBadConn
	}
	return &Stmt{
		conn:  c,
		query: query,
	}, nil
}

// Close implements driver.Conn.Close
func (c *Conn) Close() error {
	c.closed = true
	return nil
}

// Begin implements driver.Conn.Begin
func (c *Conn) Begin() (driver.Tx, error) {
	// ClickHouse doesn't support transactions in this context
	return nil, fmt.Errorf("transactions not supported")
}

// Query executes a query against the Lambda function
func (c *Conn) Query(query string) (*LambdaResponse, error) {
	if c.closed {
		return nil, driver.ErrBadConn
	}

	// Create the Lambda event payload
	event := LambdaEvent{
		RawPath: c.config.BucketPath,
		RequestContext: RequestContext{
			HTTP: HTTPContext{
				Method: "POST",
			},
		},
		Body: query,
	}

	// Marshal the event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	// Invoke the Lambda function
	input := &lambda.InvokeInput{
		FunctionName: &c.config.FunctionName,
		Payload:      eventJSON,
	}

	result, err := c.lambdaClient.Invoke(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke Lambda function: %w", err)
	}

	// Handle function errors
	if result.FunctionError != nil {
		return nil, fmt.Errorf("Lambda function error: %s", string(result.Payload))
	}

	// Parse the response
	var response LambdaResponse
	if err := json.Unmarshal(result.Payload, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Lambda function returned error: %s", response.Body)
	}

	return &response, nil
}

// ExecContext implements driver.ExecerContext
func (c *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("prepared statements with arguments not yet supported")
	}

	response, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	return &Result{
		response: response,
	}, nil
}

// QueryContext implements driver.QueryerContext
func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("prepared statements with arguments not yet supported")
	}

	response, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	return NewRows(response.Body), nil
}
