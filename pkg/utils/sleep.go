package utils

import (
	"context"
	"time"
)

// ContextSleep sleeps for the provided duration and returns
// an error if context is canceled.
func ContextSleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-timer.C:
			return nil
		}
	}
}
