package retry

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

func newBackoffPolicy() *backoff.ExponentialBackOff {
	policy := backoff.NewExponentialBackOff()
	policy.InitialInterval = 20 * time.Millisecond
	policy.MaxElapsedTime = 15 * time.Second
	return policy
}

// BackoffRetry initialize 20ms duration，last 15s
func BackoffRetry(f func() error) error {
	return backoff.Retry(f, newBackoffPolicy())
}

func BackoffRetryWithPolicy(f func() error, p *backoff.ExponentialBackOff) error {
	return backoff.Retry(f, p)
}

// Permanent wrap error, don't retry when return this kind error
func Permanent(err error) error {
	return backoff.Permanent(err)
}
