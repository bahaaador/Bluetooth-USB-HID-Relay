package retry

import (
	"math/rand"
	"time"
)

type BackoffTimer struct {
	attempts      int
	resetInterval int
	baseDelay     time.Duration
	jitterFactor  float64
}

func NewBackoffTimer(resetInterval int, baseDelay time.Duration) *BackoffTimer {
	return &BackoffTimer{
		resetInterval: resetInterval,
		baseDelay:     baseDelay,
		jitterFactor:  0.01,
	}
}

func (bt *BackoffTimer) NextDelay() time.Duration {
	// Calculate attempt number (1-based)
	currentAttempt := (bt.attempts % bt.resetInterval) + 1

	// Calculate base delay
	baseDelay := bt.baseDelay * time.Duration(currentAttempt)

	// Apply jitter: randomly adjust the delay by ±jitterFactor
	jitter := float64(baseDelay) * bt.jitterFactor * (2*rand.Float64() - 1)

	// Increment attempts counter for next time
	bt.attempts++

	return baseDelay + time.Duration(jitter)
}
