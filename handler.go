package ddb

import "github.com/aws/aws-sdk-go/service/dynamodb"

// Handler is interface for handling items from segment scan
type Handler interface {
	HandleItems(items Items) (err error)
}

// Items is the response from DDB scan
type Items []map[string]*dynamodb.AttributeValue

// HandlerFunc is a convenience type to avoid having to declare a struct
// to implement the Handler interface, it can be used like this:
//
//  scanner.Start(ddb.HandlerFunc(func(items ddb.Items) {
//    // ...
//  }))
type HandlerFunc func(items Items) (err error)

// HandleItems implements the Handler interface
func (h HandlerFunc) HandleItems(items Items) (err error) {
	err = h(items)

	return
}
