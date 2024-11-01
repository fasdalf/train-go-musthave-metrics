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
func (r *Retryer) Try(ctx context.Context, do func() error, isRetryable func(err error) bool) (r1 int, r2 error) {
	tmr := time.NewTimer(0)
	defer tmr.Stop()
out:
	for i := 0; i <= len(r.delays); i++ {
		if err := do(); err != nil {
			if !isRetryable(err) {
				r1 = i + 1
				r2 = fmt.Errorf("(retry) attempt #%d was not recoverable: %w", i+1, err)
				break
			}

			if i == len(r.delays) {
				r1 = i + 1
				r2 = fmt.Errorf("(retry) attempt #%d all attempts made, last error: %w", i+1, err)
				break
			}

			delay := r.delays[i]

			tmr.Stop()
			tmr.Reset(delay)
			select {
			case <-ctx.Done():
				r1 = i + 1
				r2 = errors.Join(fmt.Errorf("(retry) attempt #%d context exceeded after error: %w", i+1, err), context.DeadlineExceeded)
				break out
			case <-tmr.C:
				continue
			}
		}

		r1 = i + 1
		r2 = nil
		break
	}

	return r1, r2
}
