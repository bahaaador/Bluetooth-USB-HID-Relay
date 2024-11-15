package retry

import (
    "testing"
    "time"
)

func TestBackoffTimer(t *testing.T) {
    tests := []struct {
        name     string
        attempt  int
        expected time.Duration
    }{
        {
            name:     "first_attempt",
            attempt:  0,
            expected: 1 * time.Second,
        },
        {
            name:     "second_attempt",
            attempt:  1,
            expected: 2 * time.Second,
        },
        {
            name:     "third_attempt",
            attempt:  2,
            expected: 3 * time.Second,
        },
        {
            name:     "max_backoff",
            attempt:  11,
            expected: 10 * time.Second,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            timer := NewBackoffTimer()
            // Advance the timer to the correct attempt
            for i := 0; i < tt.attempt; i++ {
                timer.NextDelay()
            }
            got := timer.NextDelay()
            if got != tt.expected {
                t.Errorf("BackoffTimer attempt %d = %v, want %v",
                    tt.attempt, got, tt.expected)
            }
        })
    }
}
