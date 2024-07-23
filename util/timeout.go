package util

import (
	"context"
	"time"
)

// Util function to timeout a logic's blocking time
func TimeoutJob(logic func(), duration time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Channel to listen to the finish of logic
	logicDone := make(chan bool)

	go func() {
		logic()

		logicDone <- true
	}()

	// Race between the logic and timeout
	select {
	case <-logicDone:
		return true
	case <-ctx.Done():
		return false
	}
}
