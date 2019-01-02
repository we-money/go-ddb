package ddb

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	defaultAwsRegion     = "us-east-1"
	defaultTotalSegments = 5
)

// Config is wrapper around the configuration variables
type Config struct {
	// Svc the dynamodb connection
	Svc *dynamodb.DynamoDB

	// AwsRegion is the region the database is in. Defaults to us-west-1
	AwsRegion string

	// TableName is name of table to scan
	TableName string

	// SegmentOffset determines where to start indexing from
	SegmentOffset int

	// SegmentCount deterines how big a segment is
	SegmentCount int

	// TotalSegments determines the global amount of concurrency this will use
	TotalSegments int

	// Limit is the number of records to return during scan
	Limit int64
}

// defaults for configuration.
func (c *Config) setDefaults() {
	if c.AwsRegion == "" {
		c.AwsRegion = defaultAwsRegion
	}

	if c.TableName == "" {
		log.Fatal("TableName required as config var")
	}

	if c.TotalSegments == 0 {
		c.TotalSegments = defaultTotalSegments
	}

	if c.SegmentCount == 0 {
		c.SegmentCount = c.TotalSegments
	}

	if c.Limit == 0 {
		c.Limit = 1000
	}

	if c.Svc == nil {
		c.Svc = dynamodb.New(
			session.New(),
			aws.NewConfig().WithRegion(c.AwsRegion),
		)
	}
}
