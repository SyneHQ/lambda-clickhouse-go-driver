package clickhouse

import (
	"testing"
)

func TestParseDSN(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		want    *Config
		wantErr bool
	}{
		{
			name: "valid DSN",
			dsn:  "clickhouse-lambda://my-function@us-east-1/my-bucket/path/to/file.csv",
			want: &Config{
				FunctionName: "my-function",
				Region:       "us-east-1",
				BucketName:   "my-bucket",
				BucketPath:   "/path/to/file.csv",
			},
			wantErr: false,
		},
		{
			name: "DSN with just bucket name",
			dsn:  "clickhouse-lambda://my-function@us-east-1/my-bucket",
			want: &Config{
				FunctionName: "my-function",
				Region:       "us-east-1", 
				BucketName:   "my-bucket",
				BucketPath:   "/",
			},
			wantErr: false,
		},
		{
			name: "DSN with AWS credentials",
			dsn:  "clickhouse-lambda://my-function@us-east-1/my-bucket/path?aws_access_key_id=AKIATEST&aws_secret_access_key=secretkey&aws_session_token=token123",
			want: &Config{
				FunctionName:    "my-function",
				Region:          "us-east-1",
				BucketName:      "my-bucket",
				BucketPath:      "/path",
				AccessKeyID:     "AKIATEST",
				SecretAccessKey: "secretkey",
				SessionToken:    "token123",
			},
			wantErr: false,
		},
		{
			name: "DSN with partial AWS credentials",
			dsn:  "clickhouse-lambda://my-function@us-east-1/my-bucket?aws_access_key_id=AKIATEST&aws_secret_access_key=secretkey",
			want: &Config{
				FunctionName:    "my-function",
				Region:          "us-east-1",
				BucketName:      "my-bucket",
				BucketPath:      "/",
				AccessKeyID:     "AKIATEST",
				SecretAccessKey: "secretkey",
				SessionToken:    "",
			},
			wantErr: false,
		},
		{
			name:    "invalid scheme",
			dsn:     "mysql://user:password@localhost/db",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing function name",
			dsn:     "clickhouse-lambda://@us-east-1/my-bucket/path",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing region",
			dsn:     "clickhouse-lambda://my-function@/my-bucket/path",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL",
			dsn:     "not-a-valid-url",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDSN(tt.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.FunctionName != tt.want.FunctionName {
					t.Errorf("parseDSN() FunctionName = %v, want %v", got.FunctionName, tt.want.FunctionName)
				}
				if got.Region != tt.want.Region {
					t.Errorf("parseDSN() Region = %v, want %v", got.Region, tt.want.Region)
				}
				if got.BucketName != tt.want.BucketName {
					t.Errorf("parseDSN() BucketName = %v, want %v", got.BucketName, tt.want.BucketName)
				}
				if got.BucketPath != tt.want.BucketPath {
					t.Errorf("parseDSN() BucketPath = %v, want %v", got.BucketPath, tt.want.BucketPath)
				}
				if got.AccessKeyID != tt.want.AccessKeyID {
					t.Errorf("parseDSN() AccessKeyID = %v, want %v", got.AccessKeyID, tt.want.AccessKeyID)
				}
				if got.SecretAccessKey != tt.want.SecretAccessKey {
					t.Errorf("parseDSN() SecretAccessKey = %v, want %v", got.SecretAccessKey, tt.want.SecretAccessKey)
				}
				if got.SessionToken != tt.want.SessionToken {
					t.Errorf("parseDSN() SessionToken = %v, want %v", got.SessionToken, tt.want.SessionToken)
				}
			}
		})
	}
}

func TestNewRows(t *testing.T) {
	testData := "1\tSarah\tFox\t21\n2\tBessie\tNelson\t18\n3\tGavin\tHayes\t38"
	
	rows := NewRows(testData)
	
	// Test columns
	columns := rows.Columns()
	expectedCols := []string{"col1", "col2", "col3", "col4"}
	if len(columns) != len(expectedCols) {
		t.Errorf("Expected %d columns, got %d", len(expectedCols), len(columns))
	}
	
	for i, col := range columns {
		if col != expectedCols[i] {
			t.Errorf("Expected column %s, got %s", expectedCols[i], col)
		}
	}
	
	// Test row count
	expectedRows := 3
	if len(rows.data) != expectedRows {
		t.Errorf("Expected %d rows, got %d", expectedRows, len(rows.data))
	}
}

func TestConvertValue(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"123", int64(123)},
		{"123.45", float64(123.45)},
		{"true", true},
		{"false", false},
		{"hello", "hello"},
		{"", nil},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := convertValue(tt.input)
			if result != tt.expected {
				t.Errorf("convertValue(%q) = %v (type %T), want %v (type %T)", 
					tt.input, result, result, tt.expected, tt.expected)
			}
		})
	}
}
