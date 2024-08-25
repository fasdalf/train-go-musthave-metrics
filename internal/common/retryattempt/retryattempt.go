package retryattempt

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Retryer struct {
	delays []time.Duration
}

func NewRetryer(delays []time.Duration) *Retryer {
	return &Retryer{delays: delays}
}

func NewOneAttemptRetryer() *Retryer {
	return NewRetryer([]time.Duration{})
}

// Try Attempts retryable operation. Returns retry attempts count always and error on fail.
func (r *Retryer) Try(ctx context.Context, do func() error, isRetryable func(err error) bool) (int, error) {
	tmr := time.NewTimer(0)
	defer tmr.Stop()
	for i := 0; i <= len(r.delays); i++ {
		if err := do(); err != nil {
			if !isRetryable(err) {
				return i, fmt.Errorf("(retry) attempt #%d was not recoverable: %w", i+1, err)
			}

			if i == len(r.delays) {
				return i, fmt.Errorf("(retry) attempt #%d all attempts made, last error: %w", i+1, err)
			}

			delay := r.delays[i]

			tmr.Stop()
			tmr.Reset(delay)
			select {
			case <-ctx.Done():
				return i, errors.Join(fmt.Errorf("(retry) attempt #%d context exceeded after error: %w", i+1, err), context.DeadlineExceeded)
			case <-tmr.C:
				continue
			}
		}

		return i, nil
	}

	// unreachable
	return 0, nil
}
