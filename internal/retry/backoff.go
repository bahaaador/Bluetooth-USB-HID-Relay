package retry

import "time"

type BackoffTimer struct {
	attempts      int
	resetInterval int
}

func NewBackoffTimer() *BackoffTimer {
	return &BackoffTimer{
		resetInterval: 5,
	}
}

func (bt *BackoffTimer) NextDelay() time.Duration {
	bt.attempts++

	delay := time.Duration(bt.attempts) * time.Second

	if bt.attempts >= bt.resetInterval {
		bt.Reset()
	}

	return delay
}

func (bt *BackoffTimer) Reset() {
	bt.attempts = 0
}
