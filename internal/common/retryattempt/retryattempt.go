package retryattempt

import (
	"fmt"
	"time"
)

type Retryer struct {
	delays []time.Duration
}

func NewRetryer(delays []time.Duration) *Retryer {
	return &Retryer{delays: delays}
}

// Try Attempts retryable operation. Returns retry attempts count always and error on fail.
func (r *Retryer) Try(do func() error, isRetryable func(err error) bool) (int, error) {
	for i := 0; i <= len(r.delays); i++ {
		if err := do(); err != nil {
			if !isRetryable(err) {
				return i, fmt.Errorf("retry attempt #%d was not recoverable: %w", i+1, err)
			}

			if i == len(r.delays) {
				return i, fmt.Errorf("retry attempt #%d all attempts made, last error: %w", i+1, err)
			}

			delay := r.delays[i]
			time.Sleep(delay)
			continue
		}

		return i, nil
	}

	// unreachable
	return 0, nil
}
