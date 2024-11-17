package retry

import (
	"math"
	"testing"
	"time"
)

func TestBackoffTimer(t *testing.T) {
	tests := []struct {
		name         string
		attempt      int
		expectedBase time.Duration
	}{
		{
			name:         "first_attempt",
			attempt:      0,
			expectedBase: 1 * time.Second,
		},
		{
			name:         "second_attempt",
			attempt:      1,
			expectedBase: 2 * time.Second,
		},
		{
			name:         "third_attempt",
			attempt:      2,
			expectedBase: 3 * time.Second,
		},
		{
			name:         "fourth_attempt",
			attempt:      3,
			expectedBase: 4 * time.Second,
		},
		{
			name:         "fifth_attempt",
			attempt:      4,
			expectedBase: 5 * time.Second,
		},
		{
			name:         "reset_after_5",
			attempt:      5,
			expectedBase: 1 * time.Second,
		},
	}

	const allowedJitter = 0.2 // 20%
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := NewBackoffTimer(5, time.Second)
			// Create a new timer for each test
			for i := 0; i < tt.attempt; i++ {
				timer.NextDelay()
			}
			got := timer.NextDelay()

			difference := math.Abs(float64(got-tt.expectedBase)) / float64(tt.expectedBase)
			if difference > allowedJitter {
				t.Errorf("BackoffTimer attempt %d = %v, deviated by %.2f%% from base %v (allowed: %.0f%%)",
					tt.attempt, got, difference*100, tt.expectedBase, allowedJitter*100)
			}
		})
	}
}
