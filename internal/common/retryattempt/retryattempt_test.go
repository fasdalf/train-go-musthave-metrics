package retryattempt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
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

func TestExampleRetryer_Try(t *testing.T) {
	tests := []struct {
		name    string
		d       []time.Duration
		e       []error
		wantN   int
		wantErr bool
	}{
		{
			name:    "noop",
			d:       []time.Duration{},
			e:       []error{},
			wantN:   1,
			wantErr: false,
		},
		{
			name:    "not retried",
			d:       []time.Duration{1 * time.Nanosecond, 1 * time.Nanosecond},
			e:       []error{errors.New("0")},
			wantN:   1,
			wantErr: true,
		},
		{
			name:    "retried",
			d:       []time.Duration{10 * time.Nanosecond},
			e:       []error{io.EOF},
			wantN:   2,
			wantErr: false,
		},
		{
			name:    "timeout by errors",
			d:       []time.Duration{1 * time.Nanosecond},
			e:       []error{io.EOF, io.EOF},
			wantN:   2,
			wantErr: true,
		},
		{
			name:    "timeout by context",
			d:       []time.Duration{1 * time.Second},
			e:       []error{io.EOF, io.EOF},
			wantN:   1,
			wantErr: true,
		},
	}

	_ = NewOneAttemptRetryer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			retryer := NewRetryer(tt.d)
			i := 0
			do := func() (e error) {
				if i < len(tt.e) {
					e = tt.e[i]
				}
				i++
				return
			}
			isRetryable := func(err error) bool {
				return errors.Is(err, io.EOF)
			}
			n, err := retryer.Try(ctx, do, isRetryable)
			if n != tt.wantN {
				t.Errorf("got %d, want %d", n, tt.wantN)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
