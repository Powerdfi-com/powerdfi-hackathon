package helpers

import "time"

func retry(attempts int, sleep time.Duration, fn func() error) error {
	err := fn()
	if err != nil {
		if attempts--; attempts > 0 {
			// Exponential backoff
			time.Sleep(sleep)
			return retry(attempts, 2*sleep, fn)
		}
	}
	return err
}

func ExecuteWithRetry(fn func() error) error {
	return retry(2, 1*time.Second, fn)
}
