# ddb

DynamoDB helper for [Parallel Scan](http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/QueryAndScan.html#QueryAndScanParallelScan)

Forked from [MikeAlbertFleetSolutions/go-ddb](https://github.com/MikeAlbertFleetSolutions/go-ddb)

## Changes from original implementation

* Ability to use Index

## Installation

```
go get github.com/we-money/go-ddb
```

## Parallel Scan

To get maximum read throughput from a DynamodDB table we can leverage the [Parallel Scan](http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/QueryAndScan.html#QueryAndScanParallelScan) functionality.

```go
// structure of DDB item
type message struct {
    name string `json:"name"`
}

// new scanner with table name and total segments
scanner := ddb.NewScanner(ddb.Config{
    TableName:     "ddb-table-name",
    TotalSegments: 150,   // calculate value: (table size GB / 2GB)
    SegmentOffset: 0, // optional param for controlling offset
    SegmentCount:  150, // optional param for controlling how many routines get created
})

// start parallel scan w/ handler func
err := scanner.Start(ddb.HandlerFunc(func(items ddb.Items) (error) {
    for _, item := range items {
        var msg message
        dynamodbattribute.UnmarshalMap(item, &msg)
        fmt.Println(msg.Name)
    }
    return nil
}))

if err != nil {
  fmt.Printf("%+v\n", err)
}

```

## License

go-ddb is copyright © 2016 Clearbit. It is free software, and may
be redistributed under the terms specified in the [`LICENSE`] file.

[`LICENSE`]: /MIT-LICENSE
