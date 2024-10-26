package retryattempt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

// create and use Retryer
func ExampleRetryer_Try() {
	delays := []time.Duration{1 * time.Second, 2 * time.Second, 5 * time.Second}
	retryer := NewRetryer(delays)

	// Obtain context
	ctv := context.Background()
	// Define a function that will be retried
	do := func() error {
		// Your code here
		return nil
	}
	isRetryable := func(err error) bool {
		// Define when to retry
		return errors.Is(err, io.EOF)
	}

	if trys, err := retryer.Try(ctv, do, isRetryable); err != nil {
		fmt.Println("absolutely failed:", trys, "attempts made, ", err)
	} else {
		fmt.Println("success")
	}

	// Output:
	// success
}
