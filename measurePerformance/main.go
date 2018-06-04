package main

import (
	"context"
	"fmt"
	"time"
	"github.com/aws/aws-lambda-go/lambda"
)


func main() {
	lambda.Start(Handler)
}


func Handler(ctx context.Context, event struct{}) (string, error) {

	t := time.Now()
	ack(3, 11)
	duration := time.Since(t)

	return fmt.Sprintf("%d", duration.Nanoseconds() / 1000000), nil
}

// https://github.com/SimonWaldherr/golang-examples/blob/master/beginner/ackermann.go
func ack(n, m int64) int64 {
	for n != 0 {
		if m == 0 {
			m = 1
		} else {
			m = ack(n, m-1)
		}
		n = n - 1
	}
	return m + 1
}
