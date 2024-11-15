package retry

import "time"

type BackoffTimer struct {
    attempts   int
    maxBackoff time.Duration
}

func NewBackoffTimer() *BackoffTimer {
    return &BackoffTimer{
        maxBackoff: 10 * time.Second,
    }
}

func (bt *BackoffTimer) NextDelay() time.Duration {
    bt.attempts++
    return min(time.Duration(bt.attempts) * time.Second, bt.maxBackoff)
}

func (bt *BackoffTimer) Reset() {
    bt.attempts = 0
} 