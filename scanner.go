package ddb

import (
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jpillora/backoff"
)

// NewScanner creates a new scanner with ddb connection
func NewScanner(config Config) *Scanner {
	config.setDefaults()

	return &Scanner{
		Config: config,
	}
}

// Scanner is
type Scanner struct {
	Config
}

// Start uses the handler function to process items for each of the total shard
func (s *Scanner) Start(handler Handler) error {
	var wg sync.WaitGroup
	errored := make(chan error, 1)
	finished := make(chan bool, 1)

	for i := 0; i < s.SegmentCount; i++ {
		segment := (s.SegmentCount * s.SegmentOffset) + i

		wg.Add(1)
		go func(scan *Scanner, hand Handler, seg int) {
			defer wg.Done()

			err := scan.handlerLoop(hand, seg)
			if err != nil {
				errored <- err
			}
		}(s, handler, segment)
	}

	//	wait to be done
	go func() {
		wg.Wait()
		close(finished)
	}()

	//	handle any errors
	select {
	case <-finished:
	case goerr := <-errored:
		return goerr
	}

	return nil
}

func (s *Scanner) handlerLoop(handler Handler, segment int) (err error) {
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue

	bk := &backoff.Backoff{
		Max:    5 * time.Minute,
		Jitter: true,
	}

	for {
		// scan params
		params := &dynamodb.ScanInput{
			TableName:     aws.String(s.TableName),
			Segment:       aws.Int64(int64(segment)),
			TotalSegments: aws.Int64(int64(s.TotalSegments)),
			Limit:         aws.Int64(s.Config.Limit),
		}

		if len(s.IndexName) > 0 {
			params.IndexName = aws.String(s.IndexName)
		}

		// last evaluated key
		if lastEvaluatedKey != nil {
			params.ExclusiveStartKey = lastEvaluatedKey
		}

		// scan, sleep if rate limited
		var resp *dynamodb.ScanOutput
		resp, err = s.Svc.Scan(params)
		if err != nil {
			fmt.Println(err)
			time.Sleep(bk.Duration())
			continue
		}
		bk.Reset()

		// call the handler function with items
		err = handler.HandleItems(resp.Items)
		if err != nil {
			return
		}

		// exit if last evaluated key empty
		if resp.LastEvaluatedKey == nil {
			break
		}

		// set last evaluated key
		lastEvaluatedKey = resp.LastEvaluatedKey
	}

	return
}
